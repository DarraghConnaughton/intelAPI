package datasource

import (
	"intelligenceagent/cmd/datasource/abuseIpDb"
	"intelligenceagent/cmd/datasource/blocklistDe"
	"intelligenceagent/cmd/types"
)

func LoadDS() []types.DataSourceInterface {
	var datasources []types.DataSourceInterface
	// Load AbuseIPDb data source.
	abuseDbDS := abuseIpDb.New()
	datasources = append(datasources, &abuseDbDS)
	// Load AbuseIPDb data source.
	blockDeDS := blocklistDe.New()
	datasources = append(datasources, &blockDeDS)
	return datasources
}
