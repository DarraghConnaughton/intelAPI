package main

import (
	"log"
	"os"
	"watchdog/cmd/server"
	"watchdog/cmd/stateutil"
	"watchdog/cmd/watchdog"
)

func main() {
	log.Println("[**************************]")
	log.Println("[*] Intelligence Curator [*]")
	log.Println("[**************************]")
	var sharedState stateutil.StateManager
	errorChan := make(chan error, 1)
	go watchdog.LaunchWatchDog(errorChan, &sharedState)
	go server.StartServer(8080, errorChan, &sharedState)

	for {
		select {
		case err := <-errorChan:
			if err != nil {
				log.Println("[-]Received an error from the goroutine:", err)
				os.Exit(1)
			}
		}
	}
}
