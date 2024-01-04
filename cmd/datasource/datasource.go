package datasource

import (
	"net/http"
	"watchdog/cmd/https"
)

// IPAddressRetriever is an interface for retrieving an IP address.
type DataSource interface {
	RetrieveIPAddress() ([]string, error)
}

// DataSource represents a data source with HTTP headers, HTTPS information, and a URL.
type DataSourceConfig struct {
	Header http.Header
	HTTPS  https.HTTPS
	URL    string
	F      func(string, https.HTTPS, http.Header) ([]string, error)
}

func (dsc *DataSourceConfig) RetrieveIPAddress() ([]string, error) {
	return dsc.F(dsc.URL, dsc.HTTPS, dsc.Header)
}
