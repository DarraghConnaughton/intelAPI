package https

import (
	"io"
	"net/http"
	"strings"
)

type HTTPS struct{}

func (h *HTTPS) Get(hostname string, header http.Header) ([]byte, error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	if !strings.Contains(hostname, "https") {
		hostname = "https://" + hostname
	}

	// Set headers
	//req.Header.Set("Key", "YOUR_OWN_API_KEY")
	//req.Header
	//curl -G https://api.abuseipdb.com/api/v2/blacklist \
	//-d confidenceMinimum=90 \
	//-H "Key: YOUR_OWN_API_KEY" \
	//-H "Accept: application/json"
	//
	// Create the request
	req, err := http.NewRequest("GET", hostname, nil)
	if err != nil {
		return nil, err
	}

	req.Header = header

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
