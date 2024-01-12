package https

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

type HTTPS struct{}

func (h *HTTPS) Get(hostname string, header http.Header) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	//if strings.Contains(hostname, "http:") {
	//	hostname = strings.ReplaceAll(hostname, "http:", "https:")
	//} else if !strings.Contains(hostname, "https") {
	//	hostname = "https://" + hostname
	//}

	req, err := http.NewRequest("GET", hostname, nil)
	if err != nil {
		return nil, err
	}

	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"Unsuccessful status encountered: [%s:%d]", body, resp.StatusCode)
	}

	return body, nil
}
