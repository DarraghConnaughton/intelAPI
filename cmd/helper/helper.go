package helper

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"intelligenceagent/cmd/https"
	"log"
	"net/http"
	"sort"
	"strings"
)

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

func ExtractContentByClass(n *html.Node, targetClass string) string {
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

func SafeRetrieveBadIPAddresses(httpsClient https.HTTPS, hashes []string, urls []string) ([]string, error) {
	var ipSuperset []string
	for _, url := range urls {
		log.Println(fmt.Sprintf("[+] downloading file contents from %s", url))
		resp, err := httpsClient.Get(url, http.Header{})
		if err == nil {
			hash := md5.New()
			hash.Write(resp)
			tmp := hex.EncodeToString(hash.Sum(nil))
			if containsHash(hashes, tmp) {
				log.Println("[+] hash of downloaded contents matches expected hash, adding to blocklist.")
				ipSuperset = append(ipSuperset, strings.Split(string(resp), "\n")...)
			} else {
				log.Println("[/] warning: file hash not found in hash superset.")
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

func Prune(superset []string) []string {
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

func GetFailMode() bool {
	var tmp bool
	flag.BoolVar(&tmp, "hardfailmode", false, "whether to terminate HTTP server if goroutine produces an error.")
	return tmp
}
