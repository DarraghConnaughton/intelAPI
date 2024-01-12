package types

import (
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/stateutil"
	"net/http"
	"time"
)

type (
	AbuseIpDbData struct {
		Data []AbuseInfo `json:"data"`
	}

	BlocklistEntry struct {
		Type        string    `json:"type"`
		Updated     time.Time `json:"updated"`
		CountIPs    int       `json:"count_ips"`
		Description string    `json:"description"`
		Download    string    `json:"download"`
		MD5         string    `json:"md5"`
	}

	AbuseInfo struct {
		IPAddress            string    `json:"ipAddress"`
		CountryCode          string    `json:"countryCode,omitempty"`
		AbuseConfidenceScore int       `json:"abuseConfidenceScore,omitempty"`
		LastReportedAt       time.Time `json:"lastReportedAt,omitempty"`
	}

	DataSource struct {
		Header http.Header
		HTTPS  https.HTTPS
		URL    string
	}
)

// IPAddressRetriever is an interface for retrieving an IP address.
type DataSourceInterface interface {
	RetrieveIPAddress() ([]string, error)
	ConstructHttpHeader()
}

//type Server struct {
//	state         *stateutil.StateManager
//	routeConfiguredFlag bool
//}

// IPAddressRetriever is an interface for retrieving an IP address.
type ServerInterface interface {
	Start(int, *stateutil.StateManager)
	LoadRoutes([]RouteInfo)
}

type RouteInfo struct {
	HandlerFunc func(http.ResponseWriter, *http.Request, *stateutil.StateManager)
	Path        string
	Description string
}
