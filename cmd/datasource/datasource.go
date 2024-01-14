package datasource

import (
	"intelligenceagent/cmd/datasource/abuseIpDb"
	"intelligenceagent/cmd/datasource/blocklistDe"
	"intelligenceagent/cmd/https"
	"intelligenceagent/cmd/types"
)

func New() []types.DataSourceInterface {
	var datasources []types.DataSourceInterface
	// Load AbuseIPDb data source.
	abuseDbDS := abuseIpDb.New(https.TLSConfig{})
	datasources = append(datasources, &abuseDbDS)
	// Load AbuseIPDb data source.
	blockDeDS := blocklistDe.New(https.TLSConfig{})
	datasources = append(datasources, &blockDeDS)
	return datasources
}
