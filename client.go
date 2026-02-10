package whatsonchain

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (

	// version is the current version
	version = "v1.0.0"

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

// Client is the parent struct that contains the HTTP client.
//
// All configuration fields (apiKey, chain, network, rateLimit, userAgent)
// are stored in c.options as the single source of truth and protected by
// c.optionsMu for concurrent access. The Set* and getter methods acquire
// this mutex automatically, so callers may safely call them from any goroutine.
type Client struct {
	httpClient    HTTPInterface  // carries out the http operations
	lastRequest   *LastRequest   // is the raw information from the last request
	lastRequestMu sync.RWMutex   // protects lastRequest for concurrent access
	options       *clientOptions // single source of truth for all configuration
	optionsMu     sync.RWMutex   // protects options fields for concurrent access
}

// clientOptions holds all configuration for the client
type clientOptions struct {
	apiKey                         string
	backOffExponentFactor          float64
	backOffInitialTimeout          time.Duration
	backOffMaximumJitterInterval   time.Duration
	backOffMaxTimeout              time.Duration
	chain                          ChainType
	customHTTPClient               HTTPInterface
	dialerKeepAlive                time.Duration
	dialerTimeout                  time.Duration
	network                        NetworkType
	rateLimit                      int
	requestRetryCount              int
	requestTimeout                 time.Duration
	transportExpectContinueTimeout time.Duration
	transportIdleTimeout           time.Duration
	transportMaxIdleConnections    int
	transportTLSHandshakeTimeout   time.Duration
	userAgent                      string
}

// ClientOption is a function that modifies client options
type ClientOption func(*clientOptions)

// defaultClientOptions returns the default client options
func defaultClientOptions() *clientOptions {
	return &clientOptions{
		backOffExponentFactor:          2.0,
		backOffInitialTimeout:          2 * time.Millisecond,
		backOffMaximumJitterInterval:   2 * time.Millisecond,
		backOffMaxTimeout:              10 * time.Millisecond,
		chain:                          ChainBSV, // Default to BSV for backward compatibility
		dialerKeepAlive:                20 * time.Second,
		dialerTimeout:                  5 * time.Second,
		network:                        NetworkMain, // Default to main network
		rateLimit:                      defaultRateLimit,
		requestRetryCount:              2,
		requestTimeout:                 30 * time.Second,
		transportExpectContinueTimeout: 3 * time.Second,
		transportIdleTimeout:           20 * time.Second,
		transportMaxIdleConnections:    10,
		transportTLSHandshakeTimeout:   5 * time.Second,
		userAgent:                      defaultUserAgent,
	}
}

// WithChain sets the blockchain type (BSV or BTC)
func WithChain(chain ChainType) ClientOption {
	return func(c *clientOptions) {
		c.chain = chain
	}
}

// WithNetwork sets the network type (main, test, stn)
func WithNetwork(network NetworkType) ClientOption {
	return func(c *clientOptions) {
		c.network = network
	}
}

// WithAPIKey sets the API key for authenticated requests
func WithAPIKey(apiKey string) ClientOption {
	return func(c *clientOptions) {
		c.apiKey = apiKey
	}
}

// WithUserAgent sets a custom user agent string
func WithUserAgent(userAgent string) ClientOption {
	return func(c *clientOptions) {
		c.userAgent = userAgent
	}
}

// WithRateLimit sets the rate limit per second
func WithRateLimit(rateLimit int) ClientOption {
	return func(c *clientOptions) {
		c.rateLimit = rateLimit
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient HTTPInterface) ClientOption {
	return func(c *clientOptions) {
		c.customHTTPClient = httpClient
	}
}

// WithRequestTimeout sets the request timeout duration
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.requestTimeout = timeout
	}
}

// WithRequestRetryCount sets the number of retry attempts for failed requests
func WithRequestRetryCount(count int) ClientOption {
	return func(c *clientOptions) {
		c.requestRetryCount = count
	}
}

// WithBackoff sets the exponential backoff parameters
func WithBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.backOffInitialTimeout = initialTimeout
		c.backOffMaxTimeout = maxTimeout
		c.backOffExponentFactor = exponentFactor
		c.backOffMaximumJitterInterval = maxJitter
	}
}

// WithDialer sets the dialer configuration
func WithDialer(keepAlive, timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.dialerKeepAlive = keepAlive
		c.dialerTimeout = timeout
	}
}

// WithTransport sets the transport configuration
func WithTransport(idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int) ClientOption {
	return func(c *clientOptions) {
		c.transportIdleTimeout = idleTimeout
		c.transportTLSHandshakeTimeout = tlsTimeout
		c.transportExpectContinueTimeout = expectContinueTimeout
		c.transportMaxIdleConnections = maxIdleConnections
	}
}

// LastRequest is used to track what was submitted via the request().
// The Client protects this struct with a sync.RWMutex internally.
// The value returned by Client.LastRequest() is a copy; callers may read it freely
// but should not attempt to write back to it expecting the Client to see the change.
type LastRequest struct {
	Method     string `json:"method"`      // method is the HTTP method used
	PostData   string `json:"post_data"`   // postData is the post data submitted if POST/PUT request
	StatusCode int    `json:"status_code"` // statusCode is the last code from the request
	URL        string `json:"url"`         // url is the url used for the request
}

// newClientFromOptions creates a client from clientOptions
func newClientFromOptions(opts *clientOptions) *Client {
	// Create a client
	c := &Client{
		lastRequest: &LastRequest{},
		options:     opts,
	}

	// Is there a custom HTTP client to use?
	if opts.customHTTPClient != nil {
		c.httpClient = opts.customHTTPClient
		return c
	}

	// dial is the net dialer for clientDefaultTransport
	dial := &net.Dialer{KeepAlive: opts.dialerKeepAlive, Timeout: opts.dialerTimeout}

	// clientDefaultTransport is the default transport struct for the HTTP client
	clientDefaultTransport := &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: opts.transportExpectContinueTimeout,
		IdleConnTimeout:       opts.transportIdleTimeout,
		MaxIdleConns:          opts.transportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   opts.transportTLSHandshakeTimeout,
	}

	// Create the base HTTP client with custom transport
	baseHTTPClient := &http.Client{
		Transport: clientDefaultTransport,
		Timeout:   opts.requestTimeout,
	}

	// Determine the strategy for the http client (no retry enabled)
	if opts.requestRetryCount <= 0 {
		c.httpClient = NewSimpleHTTPClient(baseHTTPClient)
	} else { // Retry enabled
		// Create exponential back-off
		backOff := NewExponentialBackoff(
			opts.backOffInitialTimeout,
			opts.backOffMaxTimeout,
			opts.backOffExponentFactor,
			opts.backOffMaximumJitterInterval,
		)

		c.httpClient = NewRetryableHTTPClient(baseHTTPClient, opts.requestRetryCount, backOff)
	}

	return c
}
