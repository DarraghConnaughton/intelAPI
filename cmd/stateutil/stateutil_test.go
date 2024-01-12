package stateutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSecretDetection(t *testing.T) {
	stateManager := New()
	stateManager.UpdateBlocklist([]string{"1.1.1.1/32", "2.2.2.2/32"})
	assert.Equal(t, []string{"1.1.1.1/32", "2.2.2.2/32"}, stateManager.Blocklist)
	assert.Equal(t, []string{"1.1.1.1/32", "2.2.2.2/32"}, stateManager.ReadBlocklist())

	done := make(chan struct{})
	stateManager.isWriting = true
	go func() {
		// defer adds the following function to the stack, which will be executed whenever the
		// anonymous function it is embedded in ends/returns.
		defer close(done)
		stateManager.ReadBlocklist()
	}()

	select {
	// we expected an immediate return if function does not hang, which results in the closing of done.
	case <-done:
		t.Error("function returned when hanging behaviour was expected.")

		//otherwise the function is hanging (time.sleep())
	case <-time.After(1 * time.Second):
	}
}
