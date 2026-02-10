/*
Package whatsonchain is the unofficial golang implementation for the whatsonchain.com API

Example:

	// Create a new client:
	client, _ := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)

	// Get a balance for an address:
	balance, _ := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	fmt.Println("confirmed balance", balance.Confirmed)
*/
package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// EnvAPIKey is the environment variable name for the WhatsOnChain API key.
	// When set, the SDK will automatically use this key for authenticated requests
	// unless an explicit key is provided via WithAPIKey().
	EnvAPIKey = "WHATS_ON_CHAIN_API_KEY" // #nosec G101 -- env var name, not a credential
)

// NewClient creates a new WhatsOnChain client with functional options
//
// Example usage:
//
//	client, err := whatsonchain.NewClient(
//	    context.Background(),
//	    whatsonchain.WithNetwork(whatsonchain.NetworkMain),
//	    whatsonchain.WithAPIKey("your-api-key"),
//	)
func NewClient(_ context.Context, opts ...ClientOption) (ClientInterface, error) {
	// Start with defaults
	options := defaultClientOptions()

	// Apply all provided options
	for _, opt := range opts {
		opt(options)
	}

	// Auto-load API key from environment if not explicitly provided
	if options.apiKey == "" {
		if envKey := os.Getenv(EnvAPIKey); envKey != "" {
			options.apiKey = envKey
		}
	}

	// Validate chain and network
	if !validChains[options.chain] {
		return nil, ErrInvalidChain
	}
	if !validNetworks[options.network] {
		return nil, ErrInvalidNetwork
	}

	// Create and return the client
	return newClientFromOptions(options), nil
}

const (
	// maxResponseSize is the maximum response body size (50 MB).
	// Prevents unbounded memory allocation from unexpected server responses.
	maxResponseSize = 50 * 1024 * 1024
)

// request is a generic request wrapper that can be used without constraints.
// It returns the raw response body, the HTTP status code, and any error.
func (c *Client) request(ctx context.Context, url, method string, payload []byte) ([]byte, int, error) {
	// Set reader
	var bodyReader io.Reader

	// Store debugging information under mutex
	c.lastRequestMu.Lock()
	if method == http.MethodPost || method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload)
		c.lastRequest.PostData = string(payload)
	}
	c.lastRequest.Method = method
	c.lastRequest.URL = url
	c.lastRequestMu.Unlock()

	// Start the request
	var request *http.Request
	var err error
	if request, err = http.NewRequestWithContext(
		ctx, method, url, bodyReader,
	); err != nil {
		return nil, 0, err
	}

	// Read user agent and API key under options lock
	c.optionsMu.RLock()
	ua := c.options.userAgent
	apiKey := c.options.apiKey
	c.optionsMu.RUnlock()

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", ua)

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set the API key if found
	if len(apiKey) > 0 {
		request.Header.Set(apiHeaderKey, apiKey)
	}

	// Fire the http request
	var resp *http.Response
	if resp, err = c.httpClient.Do(request); err != nil {
		var statusCode int
		if resp != nil {
			statusCode = resp.StatusCode
		}
		c.lastRequestMu.Lock()
		c.lastRequest.StatusCode = statusCode
		c.lastRequestMu.Unlock()
		return nil, statusCode, err
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status under mutex
	c.lastRequestMu.Lock()
	c.lastRequest.StatusCode = resp.StatusCode
	c.lastRequestMu.Unlock()

	// Read the body with a size limit to prevent unbounded memory allocation
	var body []byte
	if body, err = io.ReadAll(io.LimitReader(resp.Body, maxResponseSize)); err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

// UserAgent will return the current user agent
func (c *Client) UserAgent() string {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.userAgent
}

// RateLimit will return the current configured rate limit
func (c *Client) RateLimit() int {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.rateLimit
}

// Chain will return the chain
func (c *Client) Chain() ChainType {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.chain
}

// Network will return the network
func (c *Client) Network() NetworkType {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.network
}

// LastRequest will return a copy of the last request information.
// The returned value is a snapshot; it is safe to read without synchronization.
func (c *Client) LastRequest() *LastRequest {
	c.lastRequestMu.RLock()
	defer c.lastRequestMu.RUnlock()
	cp := *c.lastRequest
	return &cp
}

// HTTPClient will return the current HTTP client.
// No mutex is needed here: httpClient is set once during construction
// (in newClientFromOptions) and is never modified afterward.
func (c *Client) HTTPClient() HTTPInterface {
	return c.httpClient
}

// APIKey returns the current API key
func (c *Client) APIKey() string {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.apiKey
}

// SetAPIKey sets the API key.
// This method is safe for concurrent use.
func (c *Client) SetAPIKey(apiKey string) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.apiKey = apiKey
}

// SetUserAgent sets the user agent.
// This method is safe for concurrent use.
func (c *Client) SetUserAgent(userAgent string) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.userAgent = userAgent
}

// SetRateLimit sets the rate limit.
// Values less than 1 are clamped to 1 to prevent panics in batch processors.
// This method is safe for concurrent use.
func (c *Client) SetRateLimit(rateLimit int) {
	if rateLimit < 1 {
		rateLimit = 1
	}
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.rateLimit = rateLimit
}

// SetChain sets the blockchain type.
// Returns an error if the chain type is invalid.
// This method is safe for concurrent use.
func (c *Client) SetChain(chain ChainType) error {
	if !validChains[chain] {
		return ErrInvalidChain
	}
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.chain = chain
	return nil
}

// SetNetwork sets the network type.
// Returns an error if the network type is invalid.
// This method is safe for concurrent use.
func (c *Client) SetNetwork(network NetworkType) error {
	if !validNetworks[network] {
		return ErrInvalidNetwork
	}
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.network = network
	return nil
}

// RequestTimeout returns the request timeout
func (c *Client) RequestTimeout() time.Duration {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.requestTimeout
}

// SetRequestTimeout updates the stored request timeout value.
//
// This does not rebuild the underlying HTTP client. The timeout used for
// actual requests is determined at client construction time (via WithRequestTimeout).
// This method is safe for concurrent use.
func (c *Client) SetRequestTimeout(timeout time.Duration) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.requestTimeout = timeout
}

// RequestRetryCount returns the retry count
func (c *Client) RequestRetryCount() int {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.requestRetryCount
}

// SetRequestRetryCount updates the stored retry count value.
//
// This does not rebuild the underlying HTTP client. The retry count used for
// actual requests is determined at client construction time (via WithRequestRetryCount).
// This method is safe for concurrent use.
func (c *Client) SetRequestRetryCount(count int) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.requestRetryCount = count
}

// BackoffConfig returns the backoff configuration
func (c *Client) BackoffConfig() (initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.backOffInitialTimeout, c.options.backOffMaxTimeout,
		c.options.backOffExponentFactor, c.options.backOffMaximumJitterInterval
}

// SetBackoffConfig updates the stored backoff configuration values.
//
// This does not rebuild the underlying HTTP client. The backoff configuration
// used for actual requests is determined at client construction time (via WithBackoff).
// This method is safe for concurrent use.
func (c *Client) SetBackoffConfig(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.backOffInitialTimeout = initialTimeout
	c.options.backOffMaxTimeout = maxTimeout
	c.options.backOffExponentFactor = exponentFactor
	c.options.backOffMaximumJitterInterval = maxJitter
}

// DialerConfig returns the dialer configuration
func (c *Client) DialerConfig() (keepAlive, timeout time.Duration) {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.dialerKeepAlive, c.options.dialerTimeout
}

// SetDialerConfig updates the stored dialer configuration values.
//
// This does not rebuild the underlying HTTP client. The dialer configuration
// used for actual requests is determined at client construction time (via WithDialer).
// This method is safe for concurrent use.
func (c *Client) SetDialerConfig(keepAlive, timeout time.Duration) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.dialerKeepAlive = keepAlive
	c.options.dialerTimeout = timeout
}

// TransportConfig returns the transport configuration
func (c *Client) TransportConfig() (idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int) {
	c.optionsMu.RLock()
	defer c.optionsMu.RUnlock()
	return c.options.transportIdleTimeout, c.options.transportTLSHandshakeTimeout,
		c.options.transportExpectContinueTimeout, c.options.transportMaxIdleConnections
}

// SetTransportConfig updates the stored transport configuration values.
//
// This does not rebuild the underlying HTTP client. The transport configuration
// used for actual requests is determined at client construction time (via WithTransport).
// This method is safe for concurrent use.
func (c *Client) SetTransportConfig(idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int) {
	c.optionsMu.Lock()
	defer c.optionsMu.Unlock()
	c.options.transportIdleTimeout = idleTimeout
	c.options.transportTLSHandshakeTimeout = tlsTimeout
	c.options.transportExpectContinueTimeout = expectContinueTimeout
	c.options.transportMaxIdleConnections = maxIdleConnections
}
