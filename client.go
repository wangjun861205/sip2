package sip2

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type ClientPool struct {
	length         uint64
	index          uint64
	conns          []*net.TCPConn
	host           string
	port           int
	timeout        int
	retryTimes     int
	errorDetection bool
}

func newConn(host string, port int) (*net.TCPConn, error) {
	sAddr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", sAddr)
	if err != nil {
		return nil, err
	}
	tcpConn := conn.(*net.TCPConn)
	err = tcpConn.SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}

func NewClientPool(host string, port, poolSize, timeout, retryTimes int, errorDetection bool) (*ClientPool, error) {
	conns := make([]*net.TCPConn, 0, poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := newConn(host, port)
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
	}
	return &ClientPool{uint64(poolSize), 0, conns, host, port, timeout, retryTimes, errorDetection}, nil
}

func (p *ClientPool) Pop() *net.TCPConn {
	currIndex := atomic.AddUint64(&(p.index), uint64(1))
	slicePos := (currIndex - 1) % p.length
	slotPointer := (*unsafe.Pointer)(unsafe.Pointer(&(p.conns[slicePos])))
	for {
		connPointer := atomic.SwapPointer(slotPointer, unsafe.Pointer((*net.TCPConn)(nil)))
		if connPointer == nil {
			runtime.Gosched()
			continue
		}
		conn := (*net.TCPConn)(connPointer)
		conn.SetDeadline(time.Now().Add(time.Second * time.Duration(p.timeout)))
		return conn
	}
}

func (p *ClientPool) Push(conn *net.TCPConn) {
	currIndex := atomic.AddUint64(&(p.index), uint64(1<<64-1))
	slicePos := currIndex % p.length
	slotPointer := (*unsafe.Pointer)(unsafe.Pointer(&(p.conns[slicePos])))
	for {
		ok := atomic.CompareAndSwapPointer(slotPointer, unsafe.Pointer((*net.TCPConn)(nil)), unsafe.Pointer(conn))
		if !ok {
			runtime.Gosched()
			continue
		}
		return
	}
}

func ReadResponse(conn *net.TCPConn) ([]byte, error) {
	content := make([]byte, 0, 1024)
	buffer := make([]byte, 128)
	var n int
	var err error
	for {
		n, err = conn.Read(buffer)
		if err != nil {
			return content, err
		}
		content = append(content, buffer[:n]...)
		for _, b := range buffer[:n] {
			if b == '\n' {
				return content, nil
			}
		}
	}
}

// func (p *ClientPool) ReliableCommunicate(req interface{}) (interface{}, error) {
// 	b, err := EncodeRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	conn := p.Pop()
// 	defer p.Push(conn)
// 	var bResp, zero []byte
// OUTER_WRITE:
// 	for i := 0; i < p.retryTimes; i++ {
// 		_, err = conn.Read(zero)
// 		if err != nil {
// 			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
// 				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
// 				continue OUTER_WRITE
// 			} else {
// 				var newErr error
// 				conn.Close()
// 				for {
// 					conn, newErr = newConn(p.host, p.port)
// 					if newErr == nil {
// 						conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
// 						continue OUTER_WRITE
// 					}
// 				}
// 			}
// 		}
// 		_, err = conn.Write(b)
// 		if err != nil {
// 			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
// 				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
// 				continue OUTER_WRITE
// 			} else {
// 				var newErr error
// 				conn.Close()
// 				for {
// 					conn, newErr = newConn(p.host, p.port)
// 					if newErr == nil {
// 						conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
// 						continue OUTER_WRITE
// 					}
// 				}
// 			}
// 		}
// 		break OUTER_WRITE
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// OUTER_READ:
// 	for i := 0; i < p.retryTimes; i++ {
// 		bResp, err = ReadResponse(conn)
// 		if err != nil {
// 			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
// 				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
// 				continue OUTER_READ
// 			} else {
// 				var newErr error
// 				conn.Close()
// 				for {
// 					conn, newErr = newConn(p.host, p.port)
// 					if newErr == nil {
// 						return nil, err
// 					}
// 				}
// 			}
// 		}
// 		break OUTER_READ
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := p.DecodeResponse(bResp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, ok := resp.(*ResendRequest); ok {
// 		goto OUTER_WRITE
// 	}
// 	return resp, nil
// }

func (p *ClientPool) ReliableCommunicate(req interface{}, ctx context.Context) (interface{}, error) {
	b, err := EncodeRequest(req)
	if err != nil {
		return nil, err
	}
	conn := p.Pop()
	defer p.Push(conn)
	var bResp, zero []byte
OUTER_WRITE:
	for i := 0; i < p.retryTimes; i++ {
		_, err = conn.Read(zero)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
				continue OUTER_WRITE
			} else {
				var newErr error
				conn.Close()
				for {
					conn, newErr = newConn(p.host, p.port)
					if newErr == nil {
						conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
						continue OUTER_WRITE
					}
				}
			}
		}
		_, err = conn.Write(b)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
				continue OUTER_WRITE
			} else {
				var newErr error
				conn.Close()
				for {
					conn, newErr = newConn(p.host, p.port)
					if newErr == nil {
						conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
						continue OUTER_WRITE
					}
				}
			}
		}
		break OUTER_WRITE
	}
	if err != nil {
		return nil, err
	}
OUTER_READ:
	for i := 0; i < p.retryTimes; i++ {
		bResp, err = ReadResponse(conn)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
				conn.SetDeadline(time.Now().Add(time.Duration(p.timeout) * time.Second))
				continue OUTER_READ
			} else {
				var newErr error
				conn.Close()
				for {
					conn, newErr = newConn(p.host, p.port)
					if newErr == nil {
						return nil, err
					}
				}
			}
		}
		break OUTER_READ
	}
	if err != nil {
		return nil, err
	}
	resp, err := p.DecodeResponse(bResp)
	if err != nil {
		return nil, err
	}
	if _, ok := resp.(*ResendRequest); ok {
		goto OUTER_WRITE
	}
	return resp, nil
}

func (p *ClientPool) Close(ctx context.Context) error {
	var closedNum uint64
	for {
		select {
		case <-ctx.Done():
			return errors.New("*ClientPool.Close: close error")
		default:
			p.Pop().Close()
			closedNum += 1
			if closedNum == p.length {
				return nil
			}
		}
	}
}
