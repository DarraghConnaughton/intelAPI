package abuseIpDb

import (
	"encoding/json"
	"intelligenceagent/cmd/common"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/types"
	"log"
	"net/http"
	"os"
)

// IPAddressRetriever is an interface for retrieving an IP address.
type AbuseIpDb struct {
	types.DataSourceInterface
	DataSource          types.DataSource
	ConfidenceThreshold int
	SecretLabel         string
}

func (abuseIp *AbuseIpDb) ConstructHttpHeader() {
	apiKey := os.Getenv(abuseIp.SecretLabel)
	if len(apiKey) == 0 {
		log.Println("[-] warning: no API key provided for abuseIpDb API.")
	}
	abuseIp.DataSource.Header.Set("Key", apiKey)
	abuseIp.DataSource.Header.Set("Accept", "application/json")
}

func (abuseIp *AbuseIpDb) RetrieveIPAddress() ([]string, error) {
	abuseIp.ConstructHttpHeader()
	resp, err := abuseIp.DataSource.HTTPS.Get(
		abuseIp.DataSource.URL, abuseIp.DataSource.Header)
	if err != nil {
		return nil, err
	}

	var abuseInfo types.AbuseIpDbData
	if err := json.Unmarshal(resp, &abuseInfo); err != nil {
		return nil, err
	}
	var ips []string
	for _, info := range abuseInfo.Data {
		if info.AbuseConfidenceScore > abuseIp.ConfidenceThreshold {
			ips = append(ips, info.IPAddress)
		}
	}
	return ips, nil
}

func New() AbuseIpDb {
	return AbuseIpDb{
		DataSource: types.DataSource{
			HTTPS:  https.HTTPS{},
			Header: http.Header{},
			URL:    common.AbuseIpDbApiUrl,
		},
		ConfidenceThreshold: 60,
		SecretLabel:         "ABUSEIPDB_API_TOKEN",
	}
}
