package server

import (
	"fmt"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"log"
	"net/http"
	"strings"
)

// HTTPServer extends Server and implements ServerInterface.
type HTTPServer struct {
	types.ServerInterface
	state        *stateutil.StateManager
	routeInfo    []types.RouteInfo
	BindAndServe func(string, http.Handler) error
}

func GetBlocklistHandler(w http.ResponseWriter, _ *http.Request, state *stateutil.StateManager) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join(state.ReadBlocklist(), "\n"))); err != nil || state.Mock {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) LoadRoutes(state *stateutil.StateManager) {
	for _, route := range s.routeInfo {
		localRoute := route
		http.HandleFunc(localRoute.Path, func(w http.ResponseWriter, r *http.Request) {
			localRoute.HandlerFunc(w, r, state)
		})
	}
}

func (s *HTTPServer) ListenAndServe(bind string, listenAndServe func(string, http.Handler) error) error {
	if err := listenAndServe(bind, nil); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServer) Start(port int) {
	log.Println("[+] Web Server Initialisation. Gathering routes.")
	for {
		select {
		default:
			if err := s.ListenAndServe(fmt.Sprintf(":%d", port), s.BindAndServe); err != nil {
				s.state.ErrorChan <- err
				return
			}
		}
	}
}

func New(state *stateutil.StateManager, routes []types.RouteInfo, loadRoutes bool) HTTPServer {
	httpServer := HTTPServer{
		state:        state,
		routeInfo:    routes,
		BindAndServe: http.ListenAndServe,
	}
	if loadRoutes {
		httpServer.LoadRoutes(state)
	}
	return httpServer
}
