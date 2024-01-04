package server

import (
	"net/http"
)

var PathStart = "/v1/api/"
var routes = []RouteInfo{
	{
		HandlerFunc: GetBlocklistHandler,
		//Path:        fmt.Sprintf("%s/getBlocklist", PathStart),
		Path:        "/",
		Description: "Get the current ip blocklist.",
	},
}

type RouteInfo struct {
	HandlerFunc http.HandlerFunc
	Path        string
	Description string
}

// DefineRoutes initializes and returns a slice of RouteInfo.
func RetrieveRoutes() []RouteInfo {
	return routes
}
