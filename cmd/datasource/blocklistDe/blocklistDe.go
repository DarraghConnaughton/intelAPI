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

// IPAddressRetriever is an interface for retrieving an IP address.
type BlocklistDe struct {
	types.DataSourceInterface
	DataSource  types.DataSource
	TargetClass string
}

func (bld *BlocklistDe) ConstructHttpHeader() {
	return
}

func (bld *BlocklistDe) RetrieveIPAddress() ([]string, error) {
	bld.ConstructHttpHeader()
	resp, err := bld.DataSource.HTTPS.Get(bld.DataSource.URL, bld.DataSource.Header)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(strings.NewReader(string(resp)))
	if err != nil {
		return nil, err
	}
	content := helper.ExtractContentByClass(doc, "newscontent")
	ipsuperSet, err := helper.SafeRetrieveBadIPAddresses(bld.DataSource.HTTPS,
		regexp.MustCompile(`[0-9a-fA-F]{32}`).FindAllString(content, -1),
		regexp.MustCompile(`https?://[^\s"]+\.txt`).FindAllString(content, -1))
	if err != nil {
		return nil, err
	}
	return ipsuperSet, nil
}

func New() BlocklistDe {
	return BlocklistDe{
		DataSource: types.DataSource{
			HTTPS:  https.HTTPS{},
			Header: http.Header{},
			URL:    common.BlocklistDeAPI,
		},
		TargetClass: "newscontent",
	}
}
