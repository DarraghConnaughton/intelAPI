package intelligenceagent

import (
	"fmt"
	"intelligenceagent/cmd/helper"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"log"
	"os"
	"time"
)

func MonitorErrorChannel(s *stateutil.StateManager, hardfail bool) {
	for {
		select {
		case err := <-s.ErrorChan:
			if err != nil {
				log.Println("[-]received an error from the goroutine:", err)
				if hardfail {
					log.Println("[-] hard fail mode enabled, exiting main goroutine.")
					os.Exit(1)
				}
			}
		}
	}
}

func LaunchIntelligenceAgent(datasource []types.DataSourceInterface, state *stateutil.StateManager, ticker *time.Ticker) {
	log.Println("[+] Intelligence Agent Starting [+]")
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[/] intelligence agent: tick")
				var ipSuperset []string
				for _, ds := range datasource {
					ips, err := ds.RetrieveIPAddress()
					if err != nil {
						state.ErrorChan <- fmt.Errorf("[-]Data source retrieval failed: {err: %s}", err.Error())
					}
					ipSuperset = append(ipSuperset, ips...)
				}
				state.UpdateBlocklist(helper.Prune(ipSuperset))
			}
		}
	}()
}
