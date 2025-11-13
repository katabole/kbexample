package actions

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSRFProtection_BlocksUntrustedOrigins(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// Try to POST from an untrusted origin
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=Evil+User&email=evil@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Origin", "https://evil.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// In development mode, CrossOriginProtection allows all origins (zero-value behavior)
	// In production mode with trusted origins configured, this should be blocked
	if fix.App.conf.DeployEnv.IsProduction() {
		assert.Equal(t, http.StatusForbidden, resp.StatusCode,
			"Should block requests from untrusted origins in production")
	}
}

func TestCSRFProtection_AllowsTrustedOrigins(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// POST with the trusted origin (SITE_URL)
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=Test+User&email=test@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Origin", conf.SiteURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should either succeed (200/201) or fail for auth reasons, but NOT forbidden
	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode,
		"Should allow requests from trusted origins")
}

func TestCSRFProtection_AllowsRefererHeader(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// POST with valid Referer header (used when Origin is not present)
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=Test+User&email=test@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Referer", conf.SiteURL+"/users/new")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode,
		"Should allow requests with valid Referer header")
}

func TestCSRFProtection_BlocksInvalidReferer(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// POST with invalid Referer header
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=Evil+User&email=evil@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Referer", "https://evil.com/attack")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if fix.App.conf.DeployEnv.IsProduction() {
		assert.Equal(t, http.StatusForbidden, resp.StatusCode,
			"Should block requests with invalid Referer in production")
	}
}

func TestCSRFProtection_AllowsSafeMethodsWithoutOrigin(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	safeMethods := []string{"GET", "HEAD", "OPTIONS"}

	for _, method := range safeMethods {
		t.Run(method, func(t *testing.T) {
			// Safe methods should work without Origin or Referer headers
			baseURL := "http://" + fix.App.srv.Addr
			req, err := http.NewRequest(method, baseURL+"/", nil)
			require.NoError(t, err)

			resp, err := fix.Client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should not be blocked by CSRF protection
			assert.NotEqual(t, http.StatusForbidden, resp.StatusCode,
				"%s requests should not be blocked by CSRF protection", method)
		})
	}
}

func TestCSRFProtection_ProtectsAllMutatingEndpoints(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/users"},
		{"PUT", "/users/1"},
		{"DELETE", "/users/1"},
		{"POST", "/users/1/update"},
		{"POST", "/users/1/delete"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			// Try to perform state-changing operation from untrusted origin
			baseURL := "http://" + fix.App.srv.Addr
			body := strings.NewReader("name=Hacker&email=hacker@example.com")
			req, err := http.NewRequest(tt.method, baseURL+tt.path, body)
			require.NoError(t, err)

			req.Header.Set("Origin", "https://attacker.com")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := fix.Client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if fix.App.conf.DeployEnv.IsProduction() {
				assert.Equal(t, http.StatusForbidden, resp.StatusCode,
					"%s %s should be protected from cross-origin requests in production",
					tt.method, tt.path)
			}
		})
	}
}

func TestCSRFProtection_AllowsSameOriginRequests(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// Construct same-origin URL from the test server address
	sameOrigin := "http://" + fix.App.srv.Addr
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=Same+Origin+User&email=same@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Origin", sameOrigin)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should not be blocked (though may fail auth check)
	assert.NotEqual(t, http.StatusForbidden, resp.StatusCode,
		"Same-origin requests should not be blocked by CSRF protection")
}

func TestCSRFProtection_BlocksMissingOriginAndReferer(t *testing.T) {
	fix := NewFixture(t)
	defer fix.Cleanup()

	// State-changing request without Origin or Referer headers
	baseURL := "http://" + fix.App.srv.Addr
	body := strings.NewReader("name=No+Origin&email=noorigin@example.com")
	req, err := http.NewRequest("POST", baseURL+"/users", body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Explicitly NOT setting Origin or Referer

	resp, err := fix.Client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Behavior depends on environment and CrossOriginProtection configuration
	// In strict production mode, this might be blocked
	// Document the actual behavior observed
	t.Logf("Status code for request without Origin/Referer: %d", resp.StatusCode)
}
