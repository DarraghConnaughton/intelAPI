package stateutil

import (
	"sync"
)

type StateManager struct {
	Blocklist []string
	mu        sync.Mutex
	ErrorChan chan error
	Mock      bool
}

func (sm *StateManager) UpdateBlocklist(bl []string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Blocklist = bl
}

func (sm *StateManager) ReadBlocklist() []string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return sm.Blocklist
}

func New() StateManager {
	return StateManager{
		ErrorChan: make(chan error, 1),
	}
}
