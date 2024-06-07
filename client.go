package katabole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"unicode/utf8"
)

type HTTPClientConfig struct {
	BaseURL *url.URL
}

type HTTPClient struct {
	*http.Client
	config HTTPClientConfig
}

func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		Client: http.DefaultClient,
		config: config,
	}
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.config.BaseURL == nil {
		return c.Client.Do(req)
	}

	// Clone so we don't modify the argument, just in case
	r := req.Clone(req.Context())
	r.URL.Scheme = c.config.BaseURL.Scheme
	r.URL.Host = c.config.BaseURL.Host
	r.URL.Path = path.Join(c.config.BaseURL.Path, r.URL.Path)
	r.URL.RawPath = path.Join(c.config.BaseURL.RawPath, r.URL.RawPath)
	return c.Client.Do(r)
}

// JSON
//

func (c *HTTPClient) DoJSON(req *http.Request, target any) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("got %d code and failed to read response body: %w", resp.StatusCode, err)
		}

		if !utf8.Valid(body) {
			return fmt.Errorf("got %d code and %d bytes of binary data", resp.StatusCode, len(body))
		}
		return fmt.Errorf("got %d code and response: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return fmt.Errorf("got %d code and failed to decode response body: %w", resp.StatusCode, err)
	}
	return nil
}

func (c *HTTPClient) GetJSON(urlPath string, target any) error {
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return err
	}
	return c.DoJSON(req, target)
}

func (c *HTTPClient) PutJSON(urlPath string, input any, target any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, urlPath, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return c.DoJSON(req, target)
}

func (c *HTTPClient) PostJSON(urlPath string, input any, target any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return c.DoJSON(req, target)
}

func (c *HTTPClient) DeleteJSON(urlPath string, target any) error {
	req, err := http.NewRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}
	return c.DoJSON(req, target)
}

// HTML Pages / Forms
//

func (c *HTTPClient) DoPage(req *http.Request) (string, error) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("got %d code and failed to read response body: %w", resp.StatusCode, err)
		}

		if !utf8.Valid(body) {
			return "", fmt.Errorf("got %d code and %d bytes of binary data", resp.StatusCode, len(body))
		}
		return "", fmt.Errorf("got %d code and response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("got %d code and failed to read response body: %w", resp.StatusCode, err)
	}
	return string(body), nil
}

func (c *HTTPClient) GetPage(urlPath string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return "", err
	}
	return c.DoPage(req)
}

func (c *HTTPClient) PostPage(urlPath string, input url.Values) (string, error) {
	req, err := http.NewRequest(http.MethodPost, urlPath, strings.NewReader(input.Encode()))
	if err != nil {
		return "", err
	}
	return c.DoPage(req)
}

func (c *HTTPClient) PutPage(urlPath string, input url.Values) (string, error) {
	req, err := http.NewRequest(http.MethodPut, urlPath, strings.NewReader(input.Encode()))
	if err != nil {
		return "", err
	}
	return c.DoPage(req)
}

func (c *HTTPClient) DeletePage(urlPath string) (string, error) {
	req, err := http.NewRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		return "", err
	}
	return c.DoPage(req)
}
