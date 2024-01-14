package blocklistDe

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/server"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"net/http"
	"strings"
	"testing"
	"time"
)

var MockRoutes = []types.RouteInfo{
	{
		HandlerFunc: retrieveFileWithBadHash,
		Path:        "/badFileHash.txt",
		Description: "mock endpoint whose data does not match the expected hash.",
	},
	{
		HandlerFunc: retrieveFileWithGoodHash,
		Path:        "/goodFileHash.txt",
		Description: "mock endpoint whose data matches the expected hash.",
	},
	{
		HandlerFunc: retrieveHTML,
		Path:        "/retrieveHTML",
		Description: "retrieve the HTML from blocklistde which contains other links + hash combinations",
	},
	//{
	//	HandlerFunc: retrieveBadHTML,
	//	Path:        "/retrieveBadHTML",
	//	Description: "retrieve the HTML from blocklistde which contains other links + hash combinations",
	//},
}

//func retrieveBadHTML(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
//	plainText := "This is a plain text response."
//	// Set the Content-Type header to indicate plain text
//	w.Header().Set("Content-Type", "text/plain")
//	_, err := fmt.Fprint(w, plainText)
//	if err != nil {
//		http.Error(w, "Error writing response", http.StatusInternalServerError)
//	}
//}

func retrieveHTML(w http.ResponseWriter, _ *http.Request, _ *stateutil.StateManager) {
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title class="id1">Test Page</title>
		</head>
		<body class="newscontent">
			<h1>Hello, <span>World!</span></h1>
			<p>MD5: d4d5445f37ba5f5329220cd092e106e0 
				URL: http://127.0.0.1:7777/goodFileHash.txt
            </p>
			<p>MD5: weorweorwemfmwied 
				URL: http://127.0.0.1:7777/badFileHash.txt
            </p>
		</body>
		</html>`

	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, testHTML)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func retrieveFileWithGoodHash(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join([]string{"5.5.5.5/32", "6.6.6.6/32"}, "\n"))); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func retrieveFileWithBadHash(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join([]string{"5.5.5.5/32"}, "\n"))); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func TestRetrieveIPAddress(t *testing.T) {
	blocklistde := New(https.TLSConfig{})
	blocklistde.ConstructHttpHeader()
	blocklistde.DataSource.URL = "http://127.0.0.1:7777/retrieveHTML"
	blocklistde.URLRegex = "http?://[^\\s\"]+\\.txt"

	tmpState := stateutil.New()
	s := server.New(&tmpState, MockRoutes, true)
	go s.Start(7777)
	time.Sleep(1 * time.Second)

	ips, err := blocklistde.RetrieveIPAddress()
	if err != nil {
		t.Errorf("no error expected, but one was encountered.")
	}
	assert.Equal(t, []string{"5.5.5.5/32", "6.6.6.6/32"}, ips)

	blocklistde.DataSource.URL = "http://127.0.0.1:1/retrieveBadHTML"
	ips, err = blocklistde.RetrieveIPAddress()
	if err == nil {
		t.Errorf("error expected, but none was encountered.")
	}
}
