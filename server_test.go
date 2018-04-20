package sip2

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

type ResponseHeader struct {
	Version string `json:"version"`
}

type ResponseData struct {
	Msg      string        `json:"msg"`
	Code     int           `json:"code"`
	ItemList []interface{} `json:"item_list"`
	Item     interface{}   `json:"item"`
	Meta     interface{}   `json:"meta"`
}

type JSONResponse struct {
	Header ResponseHeader `json:"header"`
	Data   ResponseData   `json:"data"`
}

func NewJSONResponse(version, msg string, code int) *JSONResponse {
	return &JSONResponse{
		Header: ResponseHeader{Version: version},
		Data: ResponseData{
			Msg:  msg,
			Code: code,
		},
	}
}

func ErrorResponse(w http.ResponseWriter, msg string, code int) {
	errResp := NewJSONResponse("2.0", msg, code)
	jsonResp, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func SuccessResponse(w http.ResponseWriter, sipResp interface{}) {
	resp := NewJSONResponse("2.0", "ok", 200)
	resp.Data.Item = sipResp
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ErrorResponse(w, "internal error", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func TestServer(t *testing.T) {
	sipServer, err := NewSIPServer("./config.json", SuccessResponse, ErrorResponse)
	if err != nil {
		log.Fatal(err)
	}
	err = sipServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
