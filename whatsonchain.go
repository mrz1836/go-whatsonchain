/*
Package whatsonchain is the unofficial golang implementation for the whatsonchain.com API

Example:

```
// Create a new client:
client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

// Get a balance for an address:
balance, _ := client.AddressBalance("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
fmt.Println("confirmed balance", balance.Confirmed)
```
*/
package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
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

	// Create and return the client
	return newClientFromOptions(options), nil
}

// request is a generic request wrapper that can be used without constraints
func (c *Client) request(ctx context.Context, url, method string, payload []byte) (response string, err error) {
	// Set reader
	var bodyReader io.Reader

	// Add post data if applicable
	if method == http.MethodPost || method == http.MethodPut {
		bodyReader = bytes.NewBuffer(payload)
		c.LastRequest().PostData = string(payload)
	}

	// Store for debugging purposes
	c.LastRequest().Method = method
	c.LastRequest().URL = url

	// Start the request
	var request *http.Request
	if request, err = http.NewRequestWithContext(
		ctx, method, url, bodyReader,
	); err != nil {
		return response, err
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
		if resp != nil {
			c.LastRequest().StatusCode = resp.StatusCode
		}
		return response, err
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Set the status
	c.LastRequest().StatusCode = resp.StatusCode

	// Read the body
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return response, err
	}

	// Return the raw JSON response
	response = string(body)
	return response, err
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

// LastRequest will return the last request information
func (c *Client) LastRequest() *LastRequest {
	return c.lastRequest
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
