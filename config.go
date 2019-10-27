package whatsonchain

import (
	"net"
	"net/http"
	"time"
)

// Package global constants and configuration
const (
	// APIEndpoint is where we fire requests
	APIEndpoint string = "https://api.whatsonchain.com/v1/bsv/"

	// NetworkMain is for main-net
	NetworkMain NetworkType = "main"

	// NetworkTest is for test-net
	NetworkTest NetworkType = "test"

	//NetworkStn is for the stn-net
	NetworkStn NetworkType = "stn"

	// ConnectionExponentFactor backoff exponent factor
	ConnectionExponentFactor float64 = 2.0

	// ConnectionInitialTimeout initial timeout
	ConnectionInitialTimeout = 2 * time.Millisecond

	// ConnectionMaximumJitterInterval jitter interval
	ConnectionMaximumJitterInterval = 2 * time.Millisecond

	// ConnectionMaxTimeout max timeout
	ConnectionMaxTimeout = 1000 * time.Millisecond

	// ConnectionRetryCount retry count
	ConnectionRetryCount int = 3

	// ConnectionWithHTTPTimeout with http timeout
	ConnectionWithHTTPTimeout = 1 * time.Second

	// ConnectionTLSHandshakeTimeout tls handshake timeout
	ConnectionTLSHandshakeTimeout = 5 * time.Second

	// ConnectionMaxIdleConnections max idle http connections
	ConnectionMaxIdleConnections int = 128

	// ConnectionIdleTimeout idle connection timeout
	ConnectionIdleTimeout = 30 * time.Second

	// ConnectionExpectContinueTimeout expect continue timeout
	ConnectionExpectContinueTimeout = 3 * time.Second

	// ConnectionDialerTimeout dialer timeout
	ConnectionDialerTimeout = 5 * time.Second

	// ConnectionDialerKeepAlive keep alive
	ConnectionDialerKeepAlive = 30 * time.Second

	// DefaultUserAgent is the default user agent for all requests
	DefaultUserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36"
)

// HTTP and Dialer connection variables
var (
	// _Dialer net dialer for ClientDefaultTransport
	_Dialer = &net.Dialer{
		KeepAlive: ConnectionDialerKeepAlive,
		Timeout:   ConnectionDialerTimeout,
	}

	// ClientDefaultTransport is the default transport struct for the HTTP client
	ClientDefaultTransport = &http.Transport{
		DialContext:           _Dialer.DialContext,
		ExpectContinueTimeout: ConnectionExpectContinueTimeout,
		IdleConnTimeout:       ConnectionIdleTimeout,
		MaxIdleConns:          ConnectionMaxIdleConnections,
		Proxy:                 http.ProxyFromEnvironment,
		TLSHandshakeTimeout:   ConnectionTLSHandshakeTimeout,
	}
)
