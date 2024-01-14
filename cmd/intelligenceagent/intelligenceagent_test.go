package intelligenceagent

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"intelligenceagent/cmd/datasource"
	"intelligenceagent/cmd/stateutil"
	"intelligenceagent/cmd/types"
	"testing"
	"time"
)

type MockDataSource struct {
	types.DataSourceInterface
	DataSource  types.DataSource
	TargetClass string
}

func (m *MockDataSource) ConstructHttpHeader() {
	return
}

func (m *MockDataSource) RetrieveIPAddress() ([]string, error) {
	if m.TargetClass == "fail" {
		return []string{}, errors.New("mock error")
	}
	return []string{"1.1.1.1/32", "2.2.2.2/32"}, nil
}

func triggerError(errChan chan error) {
	time.Sleep(1 * time.Second)
	errChan <- errors.New("trigger end of go routine")
}

func TestMonitorErrorChannel(t *testing.T) {
	state := stateutil.New()
	go triggerError(state.ErrorChan)
	// Expect this function call to hang indefinitely if called in isolation. Instead, the
	// previously launched goroutine will wake and trigger an error on a channel that the main
	// goroutine is listening to.
	if err := MonitorErrorChannel(&state, true); err == nil {
		t.Error("expected error, but encountered none.")
	}
}

func TestLaunchIntelligenceAgent(t *testing.T) {
	_ = datasource.New()
	state := stateutil.New()
	ticker := time.NewTicker(500 * time.Millisecond)
	ds := []types.DataSourceInterface{}
	ds = append(ds, &MockDataSource{
		TargetClass: "",
	})
	go LaunchIntelligenceAgent(ds, &state, ticker)
	time.Sleep(1 * time.Second)
	assert.Equal(t, []string{"1.1.1.1/32", "2.2.2.2/32"}, state.ReadBlocklist())
	ticker.Stop()
}
