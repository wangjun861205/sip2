package sip2

import (
	"context"
	"encoding/json"
	"fmt"
	"genjson"
	"net/http"
	"time"
)

// type ResponseHeader struct {
// 	Version string `json:"version"`
// }
//
// type ResponseData struct {
// 	Msg      string        `json:"msg"`
// 	Code     int           `json:"code"`
// 	ItemList []interface{} `json:"item_list"`
// 	Item     interface{}   `json:"item"`
// 	Meta     interface{}   `json:"meta"`
// }
//
// type JSONResponse struct {
// 	Header ResponseHeader `json:"header"`
// 	Data   ResponseData   `json:"data"`
// }
//
// func NewJSONResponse(version, msg string, code int) *JSONResponse {
// 	return &JSONResponse{
// 		Header: ResponseHeader{Version: version},
// 		Data: ResponseData{
// 			Msg:  msg,
// 			Code: code,
// 		},
// 	}
// }

// func ErrorResponse(w http.ResponseWriter, msg string, code int) {
// 	errResp := NewJSONResponse("2.0", msg, code)
// 	jsonResp, _ := json.Marshal(errResp)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonResp)
// }
//
// func SuccessResponse(w http.ResponseWriter, itemList []interface{}, item, meta interface{}) {
// 	resp := NewJSONResponse("2.0", "ok", 200)
// 	resp.Data.ItemList, resp.Data.Item, resp.Data.Meta = itemList, item, meta
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		ErrorResponse(w, "internal error", 500)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonResp)
// }

type SIPServer struct {
	pool     *ClientPool
	server   *http.Server
	ctx      context.Context
	cancel   context.CancelFunc
	respFunc func(http.ResponseWriter, interface{})
	errFunc  func(http.ResponseWriter, string, int)
}

func (ss *SIPServer) Route(w http.ResponseWriter, r *http.Request) {
	newCtx := context.WithValue(r.Context(), "ctx", ss.ctx)
	r = r.WithContext(newCtx)
	root := genjson.Parse(r.Body)
	if root == nil {
		ss.errFunc(w, "Not valid json format", 405)
		return
	}
	method, err := root.QueryString("header.method")
	if err != nil {
		ss.errFunc(w, "No valid method", 405)
		return
	}
	var req interface{}
	switch method {
	case "query_patron_status":
		req = NewPatronStatusRequest()
	case "query_patron_information":
		req = NewPatronInformationRequest()
	case "query_item_information":
		req = NewItemInformationRequest()
	case "check_out":
		req = NewCheckoutRequest()
	case "check_in":
		req = NewCheckinRequest()
	case "block_patron":
		req = NewBlockPatronRequest()
	case "query_sc_status":
		req = NewSCStatusRequest()
	case "login":
		req = NewLoginRequest()
	case "end_patron_session":
		req = NewEndPatronSessionRequest()
	case "fee_paid":
		req = NewFeePaidRequest()
	case "item_status_update":
		req = NewItemStatusUpdateRequest()
	case "patron_enable":
		req = NewPatronEnableRequest()
	case "hold":
		req = NewHoldRequest()
	case "renew":
		req = NewRenewRequest()
	case "renew_all":
		req = NewRenewAllRequest()
	default:
		ss.errFunc(w, "method not exist", 500)
		return
	}
	argsNode := root.Query("data")
	if argsNode == nil {
		ss.errFunc(w, "data node not exist", 500)
		return
	}
	err = json.Unmarshal([]byte(argsNode.String()), req)
	if err != nil {
		ss.errFunc(w, err.Error(), 500)
		return
	}
	resp, err := ss.pool.ReliableCommunicate(req)
	if err != nil {
		ss.errFunc(w, err.Error(), 500)
		return
	}
	// SuccessResponse(w, nil, resp, nil)
	ss.respFunc(w, resp)
}

func NewSIPServer(cfgPath string, respFunc func(http.ResponseWriter, interface{}), errFunc func(http.ResponseWriter, string, int)) (*SIPServer, error) {
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	sipServer := &SIPServer{}
	pool, err := NewClientPool(cfg.SIPConfig.Host, cfg.SIPConfig.Port, cfg.SIPConfig.PoolSize, cfg.SIPConfig.Timeout, cfg.SIPConfig.RetryTimes, cfg.SIPConfig.ErrorDetection)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", sipServer.Route)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.SIPConfig.Timeout)*time.Duration(cfg.SIPConfig.RetryTimes)*time.Second + 5*time.Second,
		WriteTimeout: time.Duration(cfg.SIPConfig.Timeout)*time.Duration(cfg.SIPConfig.RetryTimes)*time.Second + 5*time.Second,
	}
	server.SetKeepAlivesEnabled(true)
	ctx, cancel := context.WithCancel(context.Background())
	sipServer.pool = pool
	sipServer.server = server
	sipServer.ctx = ctx
	sipServer.cancel = cancel
	sipServer.respFunc = respFunc
	sipServer.errFunc = errFunc
	return sipServer, nil
}

func (ss *SIPServer) ListenAndServe() error {
	return ss.server.ListenAndServe()
}

func (ss *SIPServer) Shutdown(ctx context.Context) error {
	err := ss.pool.Close(ctx)
	err = ss.server.Shutdown(ctx)
	return err
}
