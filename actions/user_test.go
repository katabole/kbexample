package actions

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/katabole/kbexample/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUsersBasicJSON validates that we can create, get, update, and delete a user via JSON.
func TestUsersBasicJSON(t *testing.T) {
	f := NewFixture(t)
	defer f.Cleanup()

	u := models.User{ID: 1, Name: "Tim"}
	require.NoError(t, f.Client.PostJSON("/users", u, nil))

	var result models.User
	require.NoError(t, f.Client.GetJSON(fmt.Sprintf("/users/%d", u.ID), &result))
	assert.Equal(t, u, result)

	u.Name = "Tom"
	require.NoError(t, f.Client.PutJSON(fmt.Sprintf("/users/%d", u.ID), u, nil))

	require.NoError(t, f.Client.GetJSON(fmt.Sprintf("/users/%d", u.ID), &result))
	assert.Equal(t, u, result)

	require.NoError(t, f.Client.DeleteJSON(fmt.Sprintf("/users/%d", u.ID), nil))

	require.Error(t, f.Client.GetJSON(fmt.Sprintf("/users/%d", u.ID), &result))
}

// TestUsersBasicForm validates that we can create, get, update, and delete a user via html form.
func TestUsersBasicForm(t *testing.T) {
	f := NewFixture(t)
	defer f.Cleanup()

	u := models.User{ID: 1, Name: "Tim"}
	vals := url.Values{"name": []string{u.Name}}
	page, err := f.Client.PostPage("/users", vals)
	require.NoError(t, err)
	assert.Contains(t, page, u.Name)

	page, err = f.Client.GetPage(fmt.Sprintf("/users/%d", u.ID))
	require.NoError(t, err)
	assert.Contains(t, page, u.Name)

	u.Name = "Tom"
	vals = url.Values{"id": []string{"1"}, "name": []string{u.Name}}
	_, err = f.Client.PutPage(fmt.Sprintf("/users/%d", u.ID), vals)
	require.NoError(t, err)

	page, err = f.Client.GetPage(fmt.Sprintf("/users/%d", u.ID))
	require.NoError(t, err)
	assert.Contains(t, page, u.Name)

	_, err = f.Client.DeletePage(fmt.Sprintf("/users/%d", u.ID))
	require.NoError(t, err)

	_, err = f.Client.GetPage(fmt.Sprintf("/users/%d", u.ID))
	require.Error(t, err)
}
