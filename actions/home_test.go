package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHomeGET(t *testing.T) {
	f := NewFixture(t)
	defer f.Cleanup()

	page, err := f.Client.GetPage("/")
	require.NoError(t, err)
	assert.Contains(t, page, "Welcome")
}
