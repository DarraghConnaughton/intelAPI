package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"watchdog/cmd/stateutil"
)

var globalState *stateutil.StateManager

func GetBlocklistHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received.....")
	fmt.Println(globalState)
	fmt.Println(globalState.ReadBlocklist())
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(strings.Join(globalState.ReadBlocklist(), "\n"))); err != nil {
		log.Println(err)
	}
}

func loadRoutes(routeInfo []RouteInfo) {
	for _, route := range routeInfo {
		http.HandleFunc(route.Path, route.HandlerFunc)
	}
}

func StartServer(port int, errorChan chan error, state *stateutil.StateManager) {
	globalState = state
	loadRoutes(RetrieveRoutes())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		errorChan <- err
		os.Exit(1)
	}
}
