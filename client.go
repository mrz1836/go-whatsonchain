package whatsonchain

import (
	"net"
	"net/http"
	"strings"
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

// Default retry backoff. The first retry waits about 250 ms and later retries
// double the wait up to 2 s, so the default two retries span roughly 750 ms:
// long enough for a brief rate limit or a node restart to clear, short enough
// not to stall a caller. Each wait gets up to 100 ms of random jitter.
const (
	defaultBackOffInitialTimeout    = 250 * time.Millisecond
	defaultBackOffMaxTimeout        = 2 * time.Second
	defaultBackOffExponentFactor    = 2.0
	defaultBackOffMaxJitterInterval = 100 * time.Millisecond
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
	baseURL                        string
	customHTTPClient               HTTPInterface
	dialerKeepAlive                time.Duration
	dialerTimeout                  time.Duration
	network                        NetworkType
	rateLimit                      int
	requestRetryCount              int
	requestTimeout                 time.Duration
	roundTripper                   http.RoundTripper
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
		backOffExponentFactor:          defaultBackOffExponentFactor,
		backOffInitialTimeout:          defaultBackOffInitialTimeout,
		backOffMaximumJitterInterval:   defaultBackOffMaxJitterInterval,
		backOffMaxTimeout:              defaultBackOffMaxTimeout,
		baseURL:                        apiEndpointBase,
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

// validNetworks contains the set of valid network types
var validNetworks = map[NetworkType]bool{ //nolint:gochecknoglobals // read-only lookup table
	NetworkMain: true,
	NetworkTest: true,
	NetworkStn:  true,
}

// WithChain previously selected the blockchain type.
//
// Deprecated: BSV is the only supported chain. This option is now a no-op and
// is retained only for backward compatibility. It will be removed in a future
// major version.
func WithChain(_ ChainType) ClientOption {
	return func(_ *clientOptions) {}
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

// WithBaseURL overrides the base API endpoint.
//
// The default is "https://api.whatsonchain.com/v1/". The value should include
// the API version path and end with a trailing slash; a missing trailing slash
// is added automatically. The SDK still appends "bsv/<network>" and the request
// path, so a value of "https://example.com/v1/" produces requests against
// "https://example.com/v1/bsv/main/...".
//
// An empty string leaves the default unchanged. This is useful for pointing the
// client at a mirror, a reverse proxy, or a local test server (for example an
// httptest.Server) without swapping out the whole HTTP client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *clientOptions) {
		if baseURL == "" {
			return
		}
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}
		c.baseURL = baseURL
	}
}

// WithRateLimit sets the rate limit per second.
// Values less than 1 are clamped to 1 to prevent panics in batch processors.
func WithRateLimit(rateLimit int) ClientOption {
	return func(c *clientOptions) {
		if rateLimit < 1 {
			rateLimit = 1
		}
		c.rateLimit = rateLimit
	}
}

// WithHTTPClient sets a custom HTTP client, which is used as is.
//
// The SDK does not wrap the client, so its own retry logic does not run:
// WithRequestRetryCount, WithBackoff, Retry-After handling, WithRequestTimeout,
// WithDialer, WithTransport and WithRoundTripper have no effect. Any retries
// and timeouts are up to the client you pass. This option takes precedence over
// WithRoundTripper when both are set.
//
// To customize the transport (a proxy, TLS settings, instrumentation or a test
// double) and keep the SDK's retries, use WithRoundTripper instead.
func WithHTTPClient(httpClient HTTPInterface) ClientOption {
	return func(c *clientOptions) {
		c.customHTTPClient = httpClient
	}
}

// WithRoundTripper sets the http.RoundTripper that sends requests, in place of
// the SDK's default http.Transport.
//
// Unlike WithHTTPClient, the SDK still wraps it with its retries, backoff,
// Retry-After handling and request timeout. Use it to add a proxy, custom TLS
// settings, tracing or metrics, or a test double. WithDialer and WithTransport
// configure only the default transport, so they have no effect when a round
// tripper is set. A nil round tripper keeps the default transport.
// WithHTTPClient takes precedence over this option.
func WithRoundTripper(rt http.RoundTripper) ClientOption {
	return func(c *clientOptions) {
		c.roundTripper = rt
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

// WithBackoff sets the exponential backoff parameters used between retries.
//
// The wait before retry n (counting from 0) is initialTimeout * exponentFactor^n,
// capped at maxTimeout, plus a random jitter of up to maxJitter. The defaults
// are 250 ms, 2 s, 2.0 and 100 ms.
//
// On a 429 or 503 response with a Retry-After header, the client waits for the
// longer of the backoff and the Retry-After value. If Retry-After asks for more
// than maxTimeout, or the wait would outlast the request context's deadline,
// the client returns that response at once instead of retrying.
func WithBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.backOffInitialTimeout = initialTimeout
		c.backOffMaxTimeout = maxTimeout
		c.backOffExponentFactor = exponentFactor
		c.backOffMaximumJitterInterval = maxJitter
	}
}

// WithDialer sets the dialer configuration of the default transport.
// It has no effect with WithRoundTripper or WithHTTPClient.
func WithDialer(keepAlive, timeout time.Duration) ClientOption {
	return func(c *clientOptions) {
		c.dialerKeepAlive = keepAlive
		c.dialerTimeout = timeout
	}
}

// WithTransport sets the configuration of the default transport.
// It has no effect with WithRoundTripper or WithHTTPClient.
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

	// Create the base HTTP client with the custom or default transport
	baseHTTPClient := &http.Client{
		Transport: opts.roundTripper,
		Timeout:   opts.requestTimeout,
	}
	if baseHTTPClient.Transport == nil {
		baseHTTPClient.Transport = newDefaultTransport(opts)
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

// newDefaultTransport builds the http.Transport used when no round tripper is
// configured, from the dialer and transport options.
func newDefaultTransport(opts *clientOptions) *http.Transport {
	dial := &net.Dialer{KeepAlive: opts.dialerKeepAlive, Timeout: opts.dialerTimeout}

	return &http.Transport{
		DialContext:           dial.DialContext,
		ExpectContinueTimeout: opts.transportExpectContinueTimeout,
		IdleConnTimeout:       opts.transportIdleTimeout,
		MaxIdleConns:          opts.transportMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   opts.transportTLSHandshakeTimeout,
	}
}
