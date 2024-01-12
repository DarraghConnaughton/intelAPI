package server

import "intelligenceagent/cmd/types"

var Routes = []types.RouteInfo{
	{
		HandlerFunc: GetBlocklistHandler,
		Path:        "/blocklist",
		Description: "Get the current ip blocklist.",
	},
}
