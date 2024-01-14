package main

import (
	"fmt"
	"intelligenceagent/cmd/datasource"
	"intelligenceagent/cmd/helper"
	"intelligenceagent/cmd/intelligenceagent"
	s "intelligenceagent/cmd/server"
	"intelligenceagent/cmd/stateutil"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("[*] Initialising Intelligence API [*]")

	// Shared state consisting of blocklist critical section and shared error channel.
	sharedState := stateutil.New()

	// Launch Intelligence agent goroutines.
	go intelligenceagent.LaunchIntelligenceAgent(
		datasource.New(), &sharedState, time.NewTicker(10*time.Second))

	// Start HTTP server.
	server := s.New(&sharedState, s.Routes, true)
	go server.Start(8080)

	// Abort main goroutine if error is detected in one of the sub-goroutines.
	if err := intelligenceagent.MonitorErrorChannel(
		&sharedState, helper.GetFailMode()); err != nil {
		log.Println(fmt.Sprintf("[-]error encounter: %s", err.Error()))
		os.Exit(1)
	}
}
