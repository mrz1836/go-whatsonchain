package whatsonchain

import (
	"net"
	"net/http"
	"time"
)

const (

	// version is the current version
	version = "v0.13.0"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-whatsonchain: " + version

	// defaultRateLimit is the default rate limit for API requests
	defaultRateLimit int = 3

	// apiEndpointBase is where we fire requests (without chain specification)
	apiEndpointBase string = "https://api.whatsonchain.com/v1/"

	// apiHeaderKey is the header key for the API key
	apiHeaderKey string = "woc-api-key"
)

// HTTPInterface is used for the http client
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the parent struct that contains the HTTP client
type Client struct {
	apiKey      string        // optional for requests that require an API Key
	chain       ChainType     // is the blockchain type to use (BSV or BTC)
	httpClient  HTTPInterface // carries out the http operations
	lastRequest *LastRequest  // is the raw information from the last request
	network     NetworkType   // is the BitcoinSV network to use
	rateLimit   int           // configured rate limit per second
	userAgent   string        // optional for changing user agents
}

// Options holds all the configuration for connection, dialer and transport
type Options struct {
	APIKey                         string        `json:"api_key"`
	BackOffExponentFactor          float64       `json:"back_off_exponent_factor"`
	BackOffInitialTimeout          time.Duration `json:"back_off_initial_timeout"`
	BackOffMaximumJitterInterval   time.Duration `json:"back_off_maximum_jitter_interval"`
	BackOffMaxTimeout              time.Duration `json:"back_off_max_timeout"`
	DialerKeepAlive                time.Duration `json:"dialer_keep_alive"`
	DialerTimeout                  time.Duration `json:"dialer_timeout"`
	RateLimit                      int           `json:"rate_limit"`
	RequestRetryCount              int           `json:"request_retry_count"`
	RequestTimeout                 time.Duration `json:"request_timeout"`
	TransportExpectContinueTimeout time.Duration `json:"transport_expect_continue_timeout"`
	TransportIdleTimeout           time.Duration `json:"transport_idle_timeout"`
	TransportMaxIdleConnections    int           `json:"transport_max_idle_connections"`
	TransportTLSHandshakeTimeout   time.Duration `json:"transport_tls_handshake_timeout"`
	UserAgent                      string        `json:"user_agent"`
}

// LastRequest is used to track what was submitted via the request()
type LastRequest struct {
	Method     string `json:"method"`      // method is the HTTP method used
	PostData   string `json:"post_data"`   // postData is the post data submitted if POST/PUT request
	StatusCode int    `json:"status_code"` // statusCode is the last code from the request
	URL        string `json:"url"`         // url is the url used for the request
}

// ClientDefaultOptions will return an "Options" struct with the default settings
// Useful for starting with the default and then modifying as needed
func ClientDefaultOptions() (clientOptions *Options) {
	return &Options{
		BackOffExponentFactor:          2.0,
		BackOffInitialTimeout:          2 * time.Millisecond,
		BackOffMaximumJitterInterval:   2 * time.Millisecond,
		BackOffMaxTimeout:              10 * time.Millisecond,
		DialerKeepAlive:                20 * time.Second,
		DialerTimeout:                  5 * time.Second,
		RequestRetryCount:              2,
		RequestTimeout:                 30 * time.Second,
		TransportExpectContinueTimeout: 3 * time.Second,
		TransportIdleTimeout:           20 * time.Second,
		TransportMaxIdleConnections:    10,
		TransportTLSHandshakeTimeout:   5 * time.Second,
		UserAgent:                      defaultUserAgent,
		RateLimit:                      defaultRateLimit,
	}
}

// createClient will make a new http client based on the options provided
func createClient(chain ChainType, network NetworkType, options *Options, customHTTPClient HTTPInterface) (c *Client) {
	// Create a client
	c = &Client{
		chain:       chain,
		lastRequest: &LastRequest{},
		network:     network,
	}

	// Set options (either default or user modified)
	if options == nil {
		options = ClientDefaultOptions()
	}

	// Set values on the client from the given options
	c.apiKey = options.APIKey
	c.rateLimit = options.RateLimit
	c.userAgent = options.UserAgent

	// Is there a custom HTTP client to use?
	if customHTTPClient != nil {
		c.httpClient = customHTTPClient
		return c
	}

	// dial is the net dialer for clientDefaultTransport
	dial := &net.Dialer{KeepAlive: options.DialerKeepAlive, Timeout: options.DialerTimeout}

	// clientDefaultTransport is the default transport struct for the HTTP client
	clientDefaultTransport := &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: options.TransportExpectContinueTimeout,
		IdleConnTimeout:       options.TransportIdleTimeout,
		MaxIdleConns:          options.TransportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   options.TransportTLSHandshakeTimeout,
	}

	// Create the base HTTP client with custom transport
	baseHTTPClient := &http.Client{
		Transport: clientDefaultTransport,
		Timeout:   options.RequestTimeout,
	}

	// Determine the strategy for the http client (no retry enabled)
	if options.RequestRetryCount <= 0 {
		c.httpClient = NewSimpleHTTPClient(baseHTTPClient)
	} else { // Retry enabled
		// Create exponential back-off
		backOff := NewExponentialBackoff(
			options.BackOffInitialTimeout,
			options.BackOffMaxTimeout,
			options.BackOffExponentFactor,
			options.BackOffMaximumJitterInterval,
		)

		c.httpClient = NewRetryableHTTPClient(baseHTTPClient, options.RequestRetryCount, backOff)
	}

	return c
}
