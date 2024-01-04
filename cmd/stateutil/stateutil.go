package stateutil

import (
	"fmt"
	"time"
)

type StateManager struct {
	Blocklist []string
	isWriting bool
}

func (sm *StateManager) UpdateBlocklist(bl []string) {
	sm.isWriting = true
	sm.Blocklist = bl
	sm.isWriting = false
}

func (sm *StateManager) ReadBlocklist() []string {
	fmt.Println("we begin reading")
	for sm.isWriting {
		fmt.Println("are we stuck here?")
		time.Sleep(1 * time.Second)
	}
	return sm.Blocklist
}
