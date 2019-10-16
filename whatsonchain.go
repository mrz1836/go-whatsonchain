package whatsonchain

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/httpclient"
)

// NetworkType is used internally to represent the possible values
// for network in queries to be submitted: {"main", "test", "stn"}
type NetworkType string

// Client holds client configuration settings
type Client struct {

	// HTTPClient carries out the POST operations
	HTTPClient heimdall.Client

	// Parameters contains the search parameters that are submitted with your query,
	// which may affect the data returned
	Parameters *RequestParameters

	// LastRequest is the raw information from the last request
	LastRequest *LastRequest
}

// RequestParameters holds options that can affect data returned by a request.
//
// Source: https://developers.whatsonchain.com/#authentication
type RequestParameters struct {

	// UserAgent (optional for changing user agents)
	UserAgent string

	// Network is what this search should use IE: main
	Network NetworkType
}

// LastRequest is used to track what was submitted to whatsonchain on the Request()
type LastRequest struct {

	// Method is either POST or GET
	Method string

	// PostData is the post data submitted if POST request
	PostData string

	// URL is the url used for the request
	URL string
}

// NewClient creates a new search client to submit queries with.
// Parameters values are set to the defaults defined by WhatsOnChain.
//
// For more information: https://developers.whatsonchain.com/#authentication
func NewClient() (c *Client, err error) {

	// Create a client
	c = new(Client)

	// Create exponential backoff
	backOff := heimdall.NewExponentialBackoff(
		ConnectionInitialTimeout,
		ConnectionMaxTimeout,
		ConnectionExponentFactor,
		ConnectionMaximumJitterInterval,
	)

	// Create the http client
	//c.HTTPClient = new(http.Client) (@mrz this was the original HTTP client)
	c.HTTPClient = httpclient.NewClient(
		httpclient.WithHTTPTimeout(ConnectionWithHTTPTimeout),
		httpclient.WithRetrier(heimdall.NewRetrier(backOff)),
		httpclient.WithRetryCount(ConnectionRetryCount),
		httpclient.WithHTTPClient(&http.Client{
			Transport: ClientDefaultTransport,
		}),
	)

	// Create default parameters
	c.Parameters = new(RequestParameters)
	c.Parameters.UserAgent = DefaultUserAgent
	c.Parameters.Network = NetworkMain

	// Create a last request struct
	c.LastRequest = new(LastRequest)

	// Return the client
	return
}

// Request is a generic whatsonchain request wrapper that can be used without constraints
func (c *Client) Request(endpoint string, method string, params *url.Values) (response interface{}, err error) {

	// Set reader
	var bodyReader io.Reader

	// Switch on POST vs GET
	switch method {
	case "POST":
		{
			encodedParams := params.Encode()
			bodyReader = strings.NewReader(encodedParams)
			c.LastRequest.PostData = encodedParams
		}
	case "GET":
		{
			endpoint += "?" + params.Encode()
		}
	}

	// Store for debugging purposes
	c.LastRequest.Method = method
	c.LastRequest.URL = endpoint

	// Start the request
	var request *http.Request
	if request, err = http.NewRequest(method, endpoint, bodyReader); err != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", c.Parameters.UserAgent)

	// Set the content type on POST
	if method == "POST" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Fire the http request
	var resp *http.Response
	if resp, err = c.HTTPClient.Do(request); err != nil {
		return
	}

	// Close the response body
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %s", err.Error())
		}
	}()

	// Read the body
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// Parse the response
	if err = json.Unmarshal(body, response); err != nil {
		return
	}

	// Done
	return
}
