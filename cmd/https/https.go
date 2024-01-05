package https

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

type HTTPS struct{}

func (h *HTTPS) Get(hostname string, header http.Header) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if !strings.Contains(hostname, "https") {
		hostname = "https://" + hostname
	}

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
