package sip2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"genjson"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

type SIPConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	PoolSize       int    `json:"pool_size"`
	Timeout        int    `json:"timeout"`
	RetryTimes     int    `json:"retry_times"`
	ErrorDetection bool   `json:"error_detection"`
}

type ServerConfig struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	SIPConfig SIPConfig `json:"sip_config"`
}

func loadConfig(configPath string) (*ServerConfig, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bConfig, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	config := &ServerConfig{}
	err = json.Unmarshal(bConfig, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func genChecksum(bs []byte) string {
	var total int
	for _, b := range bs {
		total += int(b)
	}
	return fmt.Sprintf("%X\r", (total^0xffff)+1)
}

func checkSum(bs []byte) error {
	if len(bs) == 0 {
		return errors.New("checkSum: empty data")
	}
	content, bSum := bs[:len(bs)-5], bs[len(bs)-5:len(bs)-1]
	var s uint16
	for _, b := range content {
		s += uint16(b)
	}
	sum, err := strconv.ParseUint(string(bSum), 16, 16)
	if err != nil {
		return errors.New("checkSum: sum value not valid")
	}
	fmt.Printf("%d:%d\n", -s, sum)
	if uint16(sum) != -s {
		return errors.New("checkSum: corrupted data")
	}
	return nil
}

func formatDate() string {
	return time.Now().Format("20060102    150405")
}

func checkStructPtr(i interface{}) error {
	typ := reflect.TypeOf(i)
	kind := typ.Kind()
	if kind != reflect.Ptr {
		return errors.New("checkStructPtr: require a pointer")
	}
	if typ.Elem().Kind() != reflect.Struct {
		return errors.New("checkStructPtr: require a struct")
	}
	return nil
}

func getReqParams(r *http.Request, params ...[2]interface{}) error {
	root := genjson.Parse(r.Body)
	if root == nil {
		return errors.New("getReqParams(): not valid json format")
	}
	for _, param := range params {
		err := root.QueryValue(param[0], param[1].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func readUntil(r *bytes.Reader, delimiter byte) ([]byte, error) {
	buffer := make([]byte, 0, 1024)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == delimiter {
			break
		}
		buffer = append(buffer, b)
	}
	return buffer, nil
}

func readN(r *bytes.Reader, length int) ([]byte, error) {
	buffer := make([]byte, length)
	_, err := r.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func checkID(r *bytes.Reader, id string) error {
	idBytes := make([]byte, len(id))
	_, err := r.Read(idBytes)
	if err != nil {
		return err
	}
	if id != string(idBytes) {
		return fmt.Errorf("checkID(): field id not match (%s:%s)", id, string(idBytes))
	}
	return nil
}

func readContent(r *bytes.Reader, length int) ([]byte, error) {
	var content []byte
	var err error
	if length == -1 {
		content, err = readUntil(r, '|')
		if err != nil {
			return nil, err
		}
	} else {
		content, err = readN(r, length)
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}
