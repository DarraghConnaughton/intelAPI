package abuseIpDb

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/server"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"net/http"
	"testing"
	"time"
)

var MockRoutes = []types.RouteInfo{
	{
		HandlerFunc: retrieveAbuseDbData,
		Path:        "/abuseDB",
		Description: "retrieve abuseDBIp data for testing purposes",
	},
	{
		HandlerFunc: retrieveInvalidAbuseDbData,
		Path:        "/invalid_abuseDB",
		Description: "retrieve abuseDBIp data for testing purposes",
	},
}

func retrieveInvalidAbuseDbData(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	response, err := json.MarshalIndent([]string{"1.1.1.1/32"}, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func retrieveAbuseDbData(w http.ResponseWriter, r *http.Request, manager *stateutil.StateManager) {
	data := types.AbuseIpDbData{
		Data: []types.AbuseInfo{
			{
				IPAddress:            "192.168.1.1",
				CountryCode:          "US",
				AbuseConfidenceScore: 95,
				LastReportedAt:       time.Time{},
			},
			{
				IPAddress:            "10.0.0.1",
				CountryCode:          "CA",
				AbuseConfidenceScore: 80,
				LastReportedAt:       time.Time{},
			},
			{
				IPAddress:            "172.16.0.1",
				AbuseConfidenceScore: 70,
				LastReportedAt:       time.Time{},
			},
		},
	}
	response, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func TestRetrieveIPAddress(t *testing.T) {
	abuseipdb := New(https.TLSConfig{})
	abuseipdb.DataSource.URL = "http://127.0.0.1:9999/abuseDB"
	tmpState := stateutil.New()

	s := server.New(&tmpState, MockRoutes, true)
	go s.Start(9999)
	time.Sleep(1 * time.Second)

	ips, err := abuseipdb.RetrieveIPAddress()
	if err != nil {
		t.Errorf("no error expected, but one was encountered.")
	}
	assert.Equal(t, []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"}, ips)

	var emptySlice []string
	abuseipdb.DataSource.URL = "http://127.0.0.1:9999/invalid_abuseDB"
	ips, err = abuseipdb.RetrieveIPAddress()
	if err == nil {
		t.Errorf("no error expected, but one was encountered.")
	}
	assert.Equal(t, emptySlice, ips)
}
