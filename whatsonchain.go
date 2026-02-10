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

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", c.UserAgent())

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Set the API key if found
	if len(c.apiKey) > 0 {
		request.Header.Set(apiHeaderKey, c.apiKey)
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
	return c.userAgent
}

// RateLimit will return the current configured rate limit
func (c *Client) RateLimit() int {
	return c.rateLimit
}

// Chain will return the chain
func (c *Client) Chain() ChainType {
	return c.chain
}

// Network will return the network
func (c *Client) Network() NetworkType {
	return c.network
}

// LastRequest will return a copy of the last request information.
// The returned value is a snapshot; it is safe to read without synchronization.
func (c *Client) LastRequest() *LastRequest {
	c.lastRequestMu.RLock()
	defer c.lastRequestMu.RUnlock()
	cp := *c.lastRequest
	return &cp
}

// HTTPClient will return the current HTTP client
func (c *Client) HTTPClient() HTTPInterface {
	return c.httpClient
}

// APIKey returns the current API key
func (c *Client) APIKey() string {
	return c.apiKey
}

// SetAPIKey sets the API key
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
	if c.options != nil {
		c.options.apiKey = apiKey
	}
}

// SetUserAgent sets the user agent
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
	if c.options != nil {
		c.options.userAgent = userAgent
	}
}

// SetRateLimit sets the rate limit
func (c *Client) SetRateLimit(rateLimit int) {
	c.rateLimit = rateLimit
	if c.options != nil {
		c.options.rateLimit = rateLimit
	}
}

// SetChain sets the blockchain type
func (c *Client) SetChain(chain ChainType) {
	c.chain = chain
	if c.options != nil {
		c.options.chain = chain
	}
}

// SetNetwork sets the network type
func (c *Client) SetNetwork(network NetworkType) {
	c.network = network
	if c.options != nil {
		c.options.network = network
	}
}

// RequestTimeout returns the request timeout
func (c *Client) RequestTimeout() time.Duration {
	if c.options != nil {
		return c.options.requestTimeout
	}
	return 0
}

// SetRequestTimeout sets the request timeout
func (c *Client) SetRequestTimeout(timeout time.Duration) {
	if c.options != nil {
		c.options.requestTimeout = timeout
	}
}

// RequestRetryCount returns the retry count
func (c *Client) RequestRetryCount() int {
	if c.options != nil {
		return c.options.requestRetryCount
	}
	return 0
}

// SetRequestRetryCount sets the retry count
func (c *Client) SetRequestRetryCount(count int) {
	if c.options != nil {
		c.options.requestRetryCount = count
	}
}

// BackoffConfig returns the backoff configuration
func (c *Client) BackoffConfig() (initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) {
	if c.options != nil {
		return c.options.backOffInitialTimeout, c.options.backOffMaxTimeout,
			c.options.backOffExponentFactor, c.options.backOffMaximumJitterInterval
	}
	return 0, 0, 0, 0
}

// SetBackoffConfig sets the backoff configuration
func (c *Client) SetBackoffConfig(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) {
	if c.options != nil {
		c.options.backOffInitialTimeout = initialTimeout
		c.options.backOffMaxTimeout = maxTimeout
		c.options.backOffExponentFactor = exponentFactor
		c.options.backOffMaximumJitterInterval = maxJitter
	}
}

// DialerConfig returns the dialer configuration
func (c *Client) DialerConfig() (keepAlive, timeout time.Duration) {
	if c.options != nil {
		return c.options.dialerKeepAlive, c.options.dialerTimeout
	}
	return 0, 0
}

// SetDialerConfig sets the dialer configuration
func (c *Client) SetDialerConfig(keepAlive, timeout time.Duration) {
	if c.options != nil {
		c.options.dialerKeepAlive = keepAlive
		c.options.dialerTimeout = timeout
	}
}

// TransportConfig returns the transport configuration
func (c *Client) TransportConfig() (idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int) {
	if c.options != nil {
		return c.options.transportIdleTimeout, c.options.transportTLSHandshakeTimeout,
			c.options.transportExpectContinueTimeout, c.options.transportMaxIdleConnections
	}
	return 0, 0, 0, 0
}

// SetTransportConfig sets the transport configuration
func (c *Client) SetTransportConfig(idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int) {
	if c.options != nil {
		c.options.transportIdleTimeout = idleTimeout
		c.options.transportTLSHandshakeTimeout = tlsTimeout
		c.options.transportExpectContinueTimeout = expectContinueTimeout
		c.options.transportMaxIdleConnections = maxIdleConnections
	}
}
