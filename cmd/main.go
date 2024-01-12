package main

import (
	"intelligenceagent/cmd/datasource"
	"intelligenceagent/cmd/helper"
	"intelligenceagent/cmd/intelligenceagent"
	s "intelligenceagent/cmd/server"
	"intelligenceagent/cmd/stateutil"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("[*] Initialising Intelligence API [*]")

	// Shared state consisting of blocklist critical section and shared error channel.
	sharedState := stateutil.New()

	// Launch Intelligence agent goroutines.
	go intelligenceagent.LaunchIntelligenceAgent(
		datasource.LoadDS(), &sharedState, time.NewTicker(10*time.Second))

	// Start HTTP server.
	server := s.New(&sharedState, s.Routes, true)
	go server.Start(8080)

	// Abort main goroutine if error is detected in one of the sub-goroutines.
	intelligenceagent.MonitorErrorChannel(&sharedState, helper.GetFailMode())
}
