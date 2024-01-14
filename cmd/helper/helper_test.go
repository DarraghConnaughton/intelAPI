package helper

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/server"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

var MockRoutes = []types.RouteInfo{

	{
		HandlerFunc: retrieveFileWithBadHash,
		Path:        "/badFileHash",
		Description: "mock endpoint whose data does not match the expected hash.",
	},
	{
		HandlerFunc: retrieveFileWithGoodHash,
		Path:        "/goodFileHash",
		Description: "mock endpoint whose data matches the expected hash.",
	},
}

func retrieveFileWithBadHash(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join([]string{"5.5.5.5/32"}, "\n"))); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func retrieveFileWithGoodHash(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join([]string{"5.5.5.5/32", "6.6.6.6/32"}, "\n"))); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func TestExtractText(t *testing.T) {
	testHTML := `
		<!DOCTYPE html>
		<html>
		<head>
			<title class="id1">Test Page</title>
		</head>
		<body>
			<h1>Hello, <span>World!</span></h1>
			<p>This is a test page.</p>
		</body>
		</html>
	`

	rootNode, err := html.Parse(strings.NewReader(testHTML))
	if err != nil {
		t.Errorf("Error parsing HTML: %s", err.Error())
	}

	textContent := extractText(rootNode)
	assert.Equal(t, true, strings.Contains(textContent, "Test Page"))
	assert.Equal(t, true, strings.Contains(textContent, "Hello, World!"))
	assert.Equal(t, true, strings.Contains(textContent, "This is a test page."))

	textContent = ExtractContentByClass(rootNode, "id1")
	assert.Equal(t, textContent, "Test Page")
}

func TestSafeRetrieveBadIPAddresses(t *testing.T) {
	validHashes := []string{"93230589731e0913c32da046c2c3b40d"}
	urls := []string{"http://127.0.0.1:10000/goodFileHash", "http://127.0.0.1:10000/badFileHash"}

	tmpState := stateutil.New()
	s := server.New(&tmpState, MockRoutes, true)
	go s.Start(10000)

	time.Sleep(1 * time.Second)
	ips, err := SafeRetrieveBadIPAddresses(https.HTTPS{
		Header:    http.Header{},
		TLSConfig: https.TLSConfig{},
		Method:    "GET",
	}, validHashes, urls)
	assert.Nil(t, err)
	assert.Equal(t, []string{"5.5.5.5/32"}, ips)
}

func TestPrune(t *testing.T) {
	var nilSlice []string
	assert.Equal(t, []string{"1.1.1.1/32", "3.3.3.3/32"}, Prune([]string{"1.1.1.1/32", "1.1.1.1/32", "3.3.3.3/32"}))
	assert.Equal(t, []string{"1.1.1.1/32", "3.3.3.3/32"}, Prune([]string{"1.1.1.1/32", "3.3.3.3/32", "3.3.3.3/32"}))
	assert.Equal(t, []string{"1.1.1.1/32", "3.3.3.3/32"}, Prune([]string{"1.1.1.1/32", "3.3.3.3/32"}))
	assert.Equal(t, nilSlice, Prune([]string{}))
	assert.False(t, GetFailMode())
}
