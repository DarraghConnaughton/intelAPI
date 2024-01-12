package stateutil

import (
	"time"
)

type StateManager struct {
	Blocklist []string
	isWriting bool
	ErrorChan chan error
	Mock      bool
}

func (sm *StateManager) UpdateBlocklist(bl []string) {
	sm.isWriting = true
	sm.Blocklist = bl
	sm.isWriting = false
}

func (sm *StateManager) ReadBlocklist() []string {
	for sm.isWriting {
		time.Sleep(1 * time.Second)
	}
	return sm.Blocklist
}

func New() StateManager {
	return StateManager{
		ErrorChan: make(chan error, 1),
	}
}
