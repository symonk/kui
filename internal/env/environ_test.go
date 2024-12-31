package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigIsInEnvironment(t *testing.T) {
	t.Setenv("KUI_CONFIG", "kui.ini")
	s, ok := KeyIsInEnvironment("KUI_CONFIG")
	assert.True(t, ok)
	assert.Equal(t, "kui.ini", s)
}

func TestConfigIsNotInEnvironment(t *testing.T) {
	s, ok := KeyIsInEnvironment("KUI_CONFIG")
	assert.False(t, ok)
	assert.Empty(t, s)
}
