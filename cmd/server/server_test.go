package server

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/stateutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

//	type MockHTTPS struct {
//		types.ServerInterface
//	}
func mockListenAndServe(bind string, handler http.Handler) error {
	return errors.New("error")
}

func retrieveBlocklist(port int) []string {
	httpsClient := https.HTTPS{
		Header:    http.Header{},
		TLSConfig: https.TLSConfig{},
		Method:    "GET",
	}

	resp, _ := httpsClient.GenericMethod(
		fmt.Sprintf("http://127.0.0.1:%d/blocklist", port))
	return strings.Split(resp.Body, "\n")
}

func TestServerFunctionality(t *testing.T) {
	tmpState := stateutil.New()
	tmpState.UpdateBlocklist([]string{"3.3.3.3/32", "4.4.4.4/32"})
	server := New(&tmpState, Routes, true)

	go server.Start(8080)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, []string{"3.3.3.3/32", "4.4.4.4/32"}, retrieveBlocklist(8080))

	tmpState.Mock = true
	tmpState.UpdateBlocklist([]string{})

	// Strange output due to writer successfully writing, followed by artificial triggering of error condition.
	assert.Equal(t, []string([]string{"", ""}), retrieveBlocklist(8080))

	server2 := New(&tmpState, Routes, false)
	server2.BindAndServe = mockListenAndServe

	go server2.Start(8081)
	// Use select to wait for either an error or a timeout
	select {
	case _ = <-tmpState.ErrorChan:
		// Handle the error or check for nil if needed
		log.Println("encountered error as expected.")
	case <-time.After(100 * time.Millisecond):
		t.Error("expected error when launching server with mock ListenAndServe function")
	}
}
