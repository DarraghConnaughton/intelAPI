package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/stateutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func retrieveBlocklist(port int) []string {
	fmt.Println(fmt.Sprintf("http://127.0.0.1:%d/blocklist", port))
	httpsClient := https.HTTPS{}
	resp, _ := httpsClient.Get(
		fmt.Sprintf("http://127.0.0.1:%d/blocklist", port), http.Header{})
	return strings.Split(string(resp), "\n")
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
}
