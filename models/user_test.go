package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUsersBasic validates that we can create, get, update, and delete users.
func TestUsersBasic(t *testing.T) {
	f := NewFixture(t)
	defer f.Cleanup()

	u := &User{ID: 1, Name: "Tim"}
	newU, err := f.db.CreateUser(u)
	require.NoError(t, err)
	assert.Equal(t, u, newU)

	u.Name = "Tom"
	require.NoError(t, f.db.UpdateUser(u))
	newU, err = f.db.GetUserByID(u.ID)
	require.NoError(t, err)
	assert.Equal(t, u, newU)

	require.NoError(t, f.db.DeleteUser(u.ID))
	_, err = f.db.GetUserByID(u.ID)
	require.Error(t, err)
}
