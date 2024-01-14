package stateutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecretDetection(t *testing.T) {
	stateManager := New()
	stateManager.UpdateBlocklist([]string{"1.1.1.1/32", "2.2.2.2/32"})
	assert.Equal(t, []string{"1.1.1.1/32", "2.2.2.2/32"}, stateManager.Blocklist)
	assert.Equal(t, []string{"1.1.1.1/32", "2.2.2.2/32"}, stateManager.ReadBlocklist())
}
