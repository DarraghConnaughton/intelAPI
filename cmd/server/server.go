package server

import (
	"fmt"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"log"
	"net/http"
	"os"
	"strings"
)

// HTTPServer extends Server and implements ServerInterface.
type HTTPServer struct {
	types.ServerInterface
	state     *stateutil.StateManager
	routeInfo []types.RouteInfo
}

func GetBlocklistHandler(w http.ResponseWriter, r *http.Request, state *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	log.Println(state.ReadBlocklist())
	if _, err := w.Write([]byte(strings.Join(state.ReadBlocklist(), "\n"))); err != nil || state.Mock {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) LoadRoutes(state *stateutil.StateManager) {
	for _, route := range s.routeInfo {
		http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			route.HandlerFunc(w, r, state)
		})
	}
}

func (s *HTTPServer) Start(port int) {
	log.Println("[+] Web Server Initialisation. Gathering routes.")
	for {
		select {
		default:
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
				s.state.ErrorChan <- err
				os.Exit(1)
			}
		}
	}
}

func New(state *stateutil.StateManager, routes []types.RouteInfo, loadRoutes bool) HTTPServer {
	fmt.Println("STATE!!!!")
	fmt.Println(state)
	fmt.Println("STATE!!!!")
	httpServer := HTTPServer{
		state:     state,
		routeInfo: routes,
	}
	if loadRoutes {
		httpServer.LoadRoutes(state)
	}
	return httpServer
}
