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
)

// NewClient creates a new client for WOC requests
func NewClient(network NetworkType, clientOptions *Options, customHTTPClient HTTPInterface) ClientInterface {
	// Sets the network, options and custom HTTP client
	return createClient(network, clientOptions, customHTTPClient)
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
		return
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
		return
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
		return
	}

	// Return the raw JSON response
	response = string(body)
	return
}

// UserAgent will return the current user agent
func (c *Client) UserAgent() string {
	return c.userAgent
}

// RateLimit will return the current configured rate limit
func (c *Client) RateLimit() int {
	return c.rateLimit
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
