package blocklistDe

import (
	"golang.org/x/net/html"
	"intelligenceagent/cmd/common"
	"intelligenceagent/cmd/helper"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/types"
	"net/http"
	"regexp"
	"strings"
)

type BlocklistDe struct {
	types.DataSourceInterface
	DataSource  types.DataSource
	TargetClass string
	MD5Regex    string
	URLRegex    string
}

func (bld *BlocklistDe) ConstructHttpHeader() {
	return
}

func (bld *BlocklistDe) RetrieveIPAddress() ([]string, error) {
	resp, err := bld.DataSource.HTTPS.GenericMethod(bld.DataSource.URL)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(strings.NewReader(resp.Body))
	if err != nil {
		return nil, err
	}
	content := helper.ExtractContentByClass(doc, "newscontent")
	ipsuperSet, err := helper.SafeRetrieveBadIPAddresses(bld.DataSource.HTTPS,
		regexp.MustCompile(bld.MD5Regex).FindAllString(content, -1),
		regexp.MustCompile(bld.URLRegex).FindAllString(content, -1))
	if err != nil {
		return nil, err
	}
	return ipsuperSet, nil
}

func New(tlsconfig https.TLSConfig) BlocklistDe {
	return BlocklistDe{
		DataSource: types.DataSource{
			HTTPS: https.HTTPS{
				Header:    http.Header{},
				Method:    "GET",
				TLSConfig: tlsconfig,
			},
			URL: common.BlocklistDeAPI,
		},
		TargetClass: "newscontent",
		MD5Regex:    "[0-9a-fA-F]{32}",
		URLRegex:    "https?://[^\\s\"]+\\.txt",
	}
}
