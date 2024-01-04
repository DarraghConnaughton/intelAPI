package watchdog

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	"watchdog/cmd/common"
	"watchdog/cmd/datasource"
	"watchdog/cmd/https"
	"watchdog/cmd/stateutil"
)

type AbuseIpDb struct {
	Data []AbuseInfo `json:"data"`
}

type BlocklistEntry struct {
	Type        string    `json:"type"`
	Updated     time.Time `json:"updated"`
	CountIPs    int       `json:"count_ips"`
	Description string    `json:"description"`
	Download    string    `json:"download"`
	MD5         string    `json:"md5"`
}

type AbuseInfo struct {
	IPAddress            string    `json:"ipAddress"`
	CountryCode          string    `json:"countryCode,omitempty"`
	AbuseConfidenceScore int       `json:"abuseConfidenceScore,omitempty"`
	LastReportedAt       time.Time `json:"lastReportedAt,omitempty"`
}

// extractText extracts the text content from a given node
func extractText(n *html.Node) string {
	var textContent string
	if n.Type == html.TextNode {
		textContent = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent += extractText(c)
	}
	return textContent
}

func extractContentByClass(n *html.Node, targetClass string) string {
	var result string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && hasClass(node, targetClass) {
			result += extractText(node)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return result
}

func hasClass(n *html.Node, targetClass string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
			return true
		}
	}
	return false
}

func safeRetrieveBadIPAddresses(httpsClient https.HTTPS, hashes []string, urls []string) ([]string, error) {
	var ipSuperset []string
	for _, url := range urls {
		resp, err := httpsClient.Get(url, http.Header{})
		if err == nil {
			hash := md5.New()
			hash.Write(resp)
			if containsHash(hashes, hex.EncodeToString(hash.Sum(nil))) {
				ipSuperset = append(ipSuperset, strings.Split(string(resp), "\n")...)
			} else {
				log.Println("[**] Caution - file hash not found in hash superset.")
			}
		} else {
			return nil, err
		}
	}
	return ipSuperset, nil

}

func containsHash(hashes []string, md5hash string) bool {
	for _, hash := range hashes {
		if hash == md5hash {
			return true
		}
	}
	return false
}

func abuseIpDb(url string, httpsClient https.HTTPS, header http.Header) ([]string, error) {

	header.Set("Key", common.ApiKey)
	header.Set("Accept", "application/json")
	resp, err := httpsClient.Get(url, header)
	if err != nil {
		return nil, err
	}

	var abuseInfo AbuseIpDb
	if err := json.Unmarshal(resp, &abuseInfo); err != nil {
		return nil, err
	}

	var ips []string
	for _, info := range abuseInfo.Data {
		if info.AbuseConfidenceScore > 60 {
			ips = append(ips, info.IPAddress)
		}
	}

	return ips, nil
}

func LaunchWatchDog(errorChan chan error, state *stateutil.StateManager) {
	blDSC := datasource.DataSourceConfig{
		URL:    "https://www.blocklist.de/en/export.html",
		HTTPS:  https.HTTPS{},
		Header: http.Header{},
		F:      blocklistDe,
	}
	//abuseDSC := datasource.DataSourceConfig{
	//	URL:    apiUrl,
	//	HTTPS:  https.HTTPS{},
	//	Header: http.Header{},
	//	F:      abuseIpDb,
	//}

	// Pass ticker from main
	//tickDuration := 12 * time.Hour
	tickDuration := 10 * time.Second
	ticker := time.NewTicker(tickDuration)

	go func() {
		for {
			select {
			case <-ticker.C:
				var ipSuperset []string
				ips, err := blDSC.RetrieveIPAddress()
				if err != nil {
					errorChan <- fmt.Errorf("[-]Data source retrieval failed: {url: %s}\n[ERROR]: %s", blDSC.URL, err.Error())
					os.Exit(1)
				}
				ipSuperset = append(ipSuperset, ips...)

				//ips, err = abuseDSC.RetrieveIPAddress()
				//if err != nil {
				//	errorChan <- fmt.Errorf("[-]Data source retrieval failed: {url: %s}", abuseDSC.URL)
				//	os.Exit(1)
				//}
				ipSuperset = append(ipSuperset, ips...)
				log.Println("we make it here")
				state.UpdateBlocklist(prune(ipSuperset))
			}
		}
	}()
}

func prune(superset []string) []string {
	var finalResult []string
	sort.Slice(superset, func(i, j int) bool {
		return superset[i] < superset[j]
	})
	for i := 0; i < len(superset); i++ {
		if i == 0 || superset[i] != superset[i-1] {
			finalResult = append(finalResult, superset[i])
		}
	}
	return finalResult
}

func blocklistDe(url string, httpsClient https.HTTPS, header http.Header) ([]string, error) {
	//resp, err := httpsClient.Get("https://www.blocklist.de/en/export.html", header)
	resp, err := httpsClient.Get(url, header)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(strings.NewReader(string(resp)))
	if err != nil {
		return nil, err
	}
	// Extract content from elements with class="newscontent"
	content := extractContentByClass(doc, "newscontent")
	ipsuperSet, err := safeRetrieveBadIPAddresses(httpsClient,
		regexp.MustCompile(`[0-9a-fA-F]{32}`).FindAllString(content, -1),
		regexp.MustCompile(`https?://[^\s"]+\.txt`).FindAllString(content, -1))
	if err != nil {
		return nil, err
	}
	return ipsuperSet, nil
}
