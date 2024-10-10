package dv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	contentTypeJSON      = "application/json"
	preferFormattedValue = "odata.include-annotations=OData.Community.Display.V1.FormattedValue"
)

// NewRequest creates an [http.Request] with the provided options.
// It adds the path and query params to the base URL, and marshals the data.
func (c *Client) NewRequest(ctx context.Context, method string, path string, opts QueryOptions, data any) (*http.Request, error) {
	if path == "" {
		return nil, errors.New("new request: missing path")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	token, err := c.auth.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := c.baseURL.JoinPath(path)

	url.RawQuery = opts.ToParams().Encode()

	var body io.Reader

	if data != nil {
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		body = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating new http request: %w", err)
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	req.Header.Set("Accept", contentTypeJSON)
	req.Header.Set("Authentication", "Bearer "+token)
	req.Header.Add("Prefer", preferFormattedValue)

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

// Do makes the request using the internal [http.Client] and
// checks the error for non-200 status codes and handles it appropriately.
// Returns an [APIError] if a standard response is returned from the Dataverse server.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dataverse client network error at path '%s': %w", req.URL.String(), err)
	}

	switch {

	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return resp, nil

	case resp.StatusCode >= 400:
		return nil, handleAPIError(resp)

	default:
		return nil, fmt.Errorf("unexpected status code %d at path '%s'", resp.StatusCode, req.URL.String())

	}
}

// MakeRequest is a convenience method that calls [Client.NewRequest] then [Client.Do].
func (c *Client) MakeRequest(ctx context.Context, method string, path string, options QueryOptions, data any) (*http.Response, error) {

	req, err := c.NewRequest(ctx, method, path, options, data)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
