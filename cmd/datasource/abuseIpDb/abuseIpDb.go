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

type AbuseIpDb struct {
	types.DataSourceInterface
	ConfidenceThreshold int
	DataSource          types.DataSource
	SecretLabel         string
}

func (abuseIp *AbuseIpDb) ConstructHttpHeader() {
	apiKey := os.Getenv(abuseIp.SecretLabel)
	if len(apiKey) == 0 {
		log.Println("[-] warning: no API key provided for abuseIpDb API.")
	}
	abuseIp.DataSource.HTTPS.Header.Set("Key", apiKey)
	abuseIp.DataSource.HTTPS.Header.Set("Accept", "application/json")
}

func (abuseIp *AbuseIpDb) RetrieveIPAddress() ([]string, error) {
	resp, err := abuseIp.DataSource.HTTPS.GenericMethod(abuseIp.DataSource.URL)
	if err != nil {
		return nil, err
	}
	var abuseInfo types.AbuseIpDbData
	if err := json.Unmarshal(resp.Bytes, &abuseInfo); err != nil {
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

func New(config https.TLSConfig) AbuseIpDb {
	aidb := AbuseIpDb{
		DataSource: types.DataSource{
			HTTPS: https.HTTPS{
				Header:    http.Header{},
				Method:    "GET",
				TLSConfig: config,
			},
			URL: common.AbuseIpDbApiUrl,
		},
		ConfidenceThreshold: 60,
		SecretLabel:         "ABUSEIPDB_API_TOKEN",
	}
	aidb.ConstructHttpHeader()
	return aidb
}
