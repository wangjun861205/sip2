package sip2

import (
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	sipServer, err := NewSIPServer("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = sipServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
