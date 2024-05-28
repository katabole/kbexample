package actions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dankinder/gobase/gbexample/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersBasic(t *testing.T) {
	app, cleanup := Setup(t)
	defer cleanup()
	baseURL := "http://" + app.srv.Addr

	u := models.User{ID: 9000, Name: "Tim"}
	data, err := json.Marshal(u)
	require.Nil(t, err)
	req, err := http.NewRequest("PUT", baseURL+"/users/9000", bytes.NewReader(data))
	require.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("GET", baseURL+"/users/9000", nil)
	req.Header.Set("Accept", "application/json")
	require.Nil(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()

	var result models.User
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.Nil(t, err)
	assert.Equal(t, u, result)
}
