package katabole

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/dankinder/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mustParse(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u
}

func setup(t *testing.T) (*httpmock.MockHandlerWithHeaders, *HTTPClient, func()) {
	handler := httpmock.NewMockHandlerWithHeaders(t)
	s := httpmock.NewServer(handler)
	client := NewHTTPClient(HTTPClientConfig{BaseURL: mustParse(s.URL())})
	return handler, client, func() {
		s.Close()
		handler.AssertExpectations(t)
	}
}

func TestHTTPClientDo(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "GET", "/users/1", mock.Anything, mock.Anything).Return(httpmock.Response{
		Body: []byte(`{"name": "joebob"}`),
	})
	req, err := http.NewRequest("GET", "/users/1", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, `{"name": "joebob"}`, string(body))
}

// JSON
//

var jsonHeaderMatcher any = httpmock.HeaderMatcher("Content-Type", "application/json")

func TestHTTPClientDoJSON(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "GET", "/users/1", jsonHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`{"name": "joebob"}`),
	})

	var user struct {
		Name string
	}

	req, err := http.NewRequest("GET", "/users/1", nil)
	require.NoError(t, err)
	require.NoError(t, client.DoJSON(req, &user))
	assert.Equal(t, "joebob", user.Name)
}

type TestUser struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

func TestHTTPClientGetJSON(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "GET", "/users/1", jsonHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`{"name": "joebob"}`),
	})

	var u TestUser
	require.NoError(t, client.GetJSON("/users/1", &u))
	assert.Equal(t, "joebob", u.Name)
}

func TestHTTPClientPutJSON(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	u := &TestUser{Name: "joebob"}
	handler.On("HandleWithHeaders", "PUT", "/users/9000", jsonHeaderMatcher, httpmock.JSONMatcher(u)).Return(httpmock.Response{
		Body: []byte(`{"id": 9000, "name": "joebob"}`),
	})

	var payload TestUser
	require.NoError(t, client.PutJSON("/users/9000", &u, &payload))
	assert.Equal(t, 9000, payload.ID)
	assert.Equal(t, "joebob", payload.Name)
}

func TestHTTPClientPutJSONNoTarget(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()
	u := &TestUser{Name: "joebob"}
	handler.On("HandleWithHeaders", "PUT", "/users/9000", jsonHeaderMatcher, httpmock.JSONMatcher(u)).Return(httpmock.Response{
		Body: []byte(`{"id": 9000, "name": "joebob"}`),
	})

	require.NoError(t, client.PutJSON("/users/9000", &u, nil))
}

func TestHTTPClientPostJSON(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	u := &TestUser{Name: "joebob"}
	handler.On("HandleWithHeaders", "POST", "/users", jsonHeaderMatcher, httpmock.JSONMatcher(u)).Return(httpmock.Response{
		Body: []byte(`{"id": 9000, "name": "joebob"}`),
	})

	var payload TestUser
	require.NoError(t, client.PostJSON("/users", &u, &payload))
	assert.Equal(t, 9000, payload.ID)
	assert.Equal(t, "joebob", payload.Name)
}

func TestHTTPClientPostJSONNoTarget(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	u := &TestUser{Name: "joebob"}
	handler.On("HandleWithHeaders", "POST", "/users", jsonHeaderMatcher, httpmock.JSONMatcher(u)).Return(httpmock.Response{
		Body: []byte(`{"id": 9000, "name": "joebob"}`),
	})

	require.NoError(t, client.PostJSON("/users", &u, nil))
}

func TestHTTPClientDeleteJSON(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "DELETE", "/users/9000", jsonHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`{"result": "ok"}`),
	})

	payload := map[string]string{}
	require.NoError(t, client.DeleteJSON("/users/9000", &payload))
	assert.Equal(t, "ok", payload["result"])
}

func TestHTTPClientDeleteJSONNoTarget(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "DELETE", "/users/9000", jsonHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`{"result": "ok"}`),
	})

	require.NoError(t, client.DeleteJSON("/users/9000", nil))
}

// HTML Pages / Forms
//

var htmlHeaderMatcher any = httpmock.HeaderMatcher("Content-Type", "application/x-www-form-urlencoded")

func TestHTTPClientDoPage(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "GET", "/users/1", htmlHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`<html>hi there</html>`),
	})

	req, err := http.NewRequest("GET", "/users/1", nil)
	require.NoError(t, err)
	page, err := client.DoPage(req)
	require.NoError(t, err)
	assert.Contains(t, page, "hi there")
}

func TestHTTPClientGetPage(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "GET", "/users/1", htmlHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`<html>hi there</html>`),
	})

	page, err := client.GetPage("/users/1")
	require.NoError(t, err)
	assert.Contains(t, page, "hi there")
}

func TestHTTPClientPostPage(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	vals := url.Values{"name": []string{"joebob"}}
	handler.On("HandleWithHeaders", "POST", "/users", htmlHeaderMatcher, []byte(vals.Encode())).Return(httpmock.Response{
		Body: []byte(`<html>hi there</html>`),
	})

	page, err := client.PostPage("/users", vals)
	require.NoError(t, err)
	assert.Contains(t, page, "hi there")
}

func TestHTTPClientPutPage(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	vals := url.Values{"name": []string{"joebob"}}
	handler.On("HandleWithHeaders", "PUT", "/users/9000", htmlHeaderMatcher, []byte(vals.Encode())).Return(httpmock.Response{
		Body: []byte(`<html>hi there</html>`),
	})

	page, err := client.PutPage("/users/9000", vals)
	require.NoError(t, err)
	assert.Contains(t, page, "hi there")
}

func TestHTTPClientDeletePage(t *testing.T) {
	handler, client, cleanup := setup(t)
	defer cleanup()

	handler.On("HandleWithHeaders", "DELETE", "/users/9000", htmlHeaderMatcher, mock.Anything).Return(httpmock.Response{
		Body: []byte(`<html>hi there</html>`),
	})

	page, err := client.DeletePage("/users/9000")
	require.NoError(t, err)
	assert.Contains(t, page, "hi there")
}
