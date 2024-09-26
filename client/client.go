package dataverse

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Version = "1.00"

	defaultUserAgent  = "go-dataverse" + "/" + Version
	defaultHTTPClient = &http.Client{Timeout: time.Second * 15}

	ErrMissingConfig = errors.New("missing config")
)

// TokenGetter represents an auth client that gets a token.
type TokenGetter interface {
	GetToken(context.Context) (string, error)
}

// Config holds the necessary configuration strings.
type Config struct {
	BaseURL      string
	TenantID     string
	ClientID     string
	ClientSecret string
}

// Validate returns an error that lists the missing fields.
func (c Config) Validate() error {
	errors := []string{}
	switch {
	case c.BaseURL == "":
		errors = append(errors, "BaseURL")
		fallthrough
	case c.TenantID == "":
		errors = append(errors, "TenantID")
		fallthrough
	case c.ClientID == "":
		errors = append(errors, "ClientID")
		fallthrough
	case c.ClientSecret == "":
		errors = append(errors, "ClientSecret")
	}
	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingConfig, strings.Join(errors, ", "))
	}
	return nil
}

// Client interacts with the Dataverse API.
// It handles authentication, decoding JSON, and errors.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	auth       TokenGetter
	userAgent  string
}

// ClientOption is an optional value that modifies the client.
type ClientOption func(*Client)

// NewClient returns an instance of Client. It creates an [AuthClient] internally that implements [TokenGetter].
// It takes a variadic set of [ClientOption].
// Current provided options are [WithHTTPClient], [WithUserAgent], [WithAuthClient].
func NewClient(config Config, options ...ClientOption) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing BaseURL: %w", err)
	}

	client := &Client{baseURL: parsedURL}

	for _, opt := range options {
		opt(client)
	}

	// Set defaults
	if client.auth == nil {
		authClient, err := NewSecretAuthClient(config.TenantID, config.ClientID, config.ClientSecret)
		if err != nil {
			return nil, err
		}
		client.auth = authClient
	}

	client.userAgent = cmp.Or(client.userAgent, defaultUserAgent)
	client.httpClient = cmp.Or(client.httpClient, defaultHTTPClient)

	return client, nil
}

// WithHTTPClient sets the underlying [*http.Client].
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithUserAgent sets the user agent.
func WithUserAgent(agent string) ClientOption {
	return func(c *Client) {
		c.userAgent = agent
	}
}

// WithAuthClient sets the auth client that implements [TokenGetter].
func WithAuthClient(client TokenGetter) ClientOption {
	return func(c *Client) {
		c.auth = client
	}
}
