package server

import (
	"net/http"
)

var PathStart = "/v1/api/"
var routes = []RouteInfo{
	{
		HandlerFunc: GetBlocklistHandler,
		//Path:        fmt.Sprintf("%s/getBlocklist", PathStart),
		Path:        "/blocklist",
		Description: "Get the current ip blocklist.",
	},
}

type RouteInfo struct {
	HandlerFunc http.HandlerFunc
	Path        string
	Description string
}

func RetrieveRoutes() []RouteInfo {
	return routes
}
