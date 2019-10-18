package whatsonchain

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// StatusCode is the last code from the request
	StatusCode int

	// URL is the url used for the request
	URL string
}

// NewClient creates a new client to submit queries with.
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
func (c *Client) Request(endpoint string, method string, params *url.Values, payload []byte) (response string, err error) {

	// Set reader
	var bodyReader io.Reader

	// Add the network value
	endpoint = fmt.Sprintf("%s%s/%s", APIEndpoint, c.Parameters.Network, endpoint)

	// Switch on POST vs GET
	switch method {
	case "POST":
		{
			if len(payload) > 0 {
				bodyReader = bytes.NewBuffer(payload)
			} else {
				encodedParams := params.Encode()
				bodyReader = strings.NewReader(encodedParams)
				c.LastRequest.PostData = encodedParams
			}
		}
	case "GET":
		{
			if params != nil {
				endpoint += "?" + params.Encode()
			}
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
		if len(payload) > 0 {
			request.Header.Set("Content-Type", "application/json")
		} else {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Fire the http request
	var resp *http.Response
	if resp, err = c.HTTPClient.Do(request); err != nil {
		return
	}

	// Close the response body
	defer func() {
		if bodyErr := resp.Body.Close(); bodyErr != nil {
			log.Printf("error closing response body: %s", bodyErr.Error())
		}
	}()

	// Save the status
	c.LastRequest.StatusCode = resp.StatusCode

	// Read the body
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// Parse the response
	response = string(body)

	// Done
	return
}

// GetHealth gets the status from whatsonchain
//
// For more information: https://developers.whatsonchain.com/#health
func (c *Client) GetHealth() (status string, err error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/woc

	return c.Request("woc", "GET", nil, nil)
}

// GetChainInfo gets the chain info from whatsonchain
//
// For more information: https://developers.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo() (chainInfo *ChainInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/info
	resp, err = c.Request("chain/info", "GET", nil, nil)
	if err != nil {
		return
	}

	chainInfo = new(ChainInfo)
	if err = json.Unmarshal([]byte(resp), chainInfo); err != nil {
		return
	}
	return
}

// GetBlockByHash gets the block info
//
// Form more information: https://developers.whatsonchain.com/#get-by-hash
func (c *Client) GetBlockByHash(hash string) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>
	resp, err = c.Request("block/hash/"+hash, "GET", nil, nil)
	if err != nil {
		return
	}

	blockInfo = new(BlockInfo)
	if err = json.Unmarshal([]byte(resp), blockInfo); err != nil {
		return
	}
	return

}

// GetBlockByHeight gets the block info
//
// Form more information: https://developers.whatsonchain.com/#get-by-height
func (c *Client) GetBlockByHeight(height int64) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/height/<height>
	resp, err = c.Request(fmt.Sprintf("block/height/%d", height), "GET", nil, nil)
	if err != nil {
		return
	}

	blockInfo = new(BlockInfo)
	if err = json.Unmarshal([]byte(resp), blockInfo); err != nil {
		return
	}
	return

}

// GetBlockPages gets the block info
//
// Form more information: https://developers.whatsonchain.com/#get-block-pages
func (c *Client) GetBlockPages(hash string, page int) (txList *BlockPagesInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>/page/1
	resp, err = c.Request(fmt.Sprintf("block/hash/%s/page/%d", hash, page), "GET", nil, nil)
	if err != nil {
		return
	}

	txList = new(BlockPagesInfo)
	if err = json.Unmarshal([]byte(resp), txList); err != nil {
		return
	}
	return

}

// GetTxByHash gets the tx info
//
// Form more information: https://developers.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(hash string) (txInfo *TxInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/hash/<hash>
	resp, err = c.Request("tx/hash/"+hash, "GET", nil, nil)
	if err != nil {
		return
	}

	txInfo = new(TxInfo)
	if err = json.Unmarshal([]byte(resp), txInfo); err != nil {
		return
	}
	return

}

// BroadcastTx will broadcast a transaction via whatsonchain
//
// Form more information: https://developers.whatsonchain.com/#broadcast-transaction
func (c *Client) BroadcastTx(txHex string) (txID string, err error) {

	// Start the post data
	stringVal := fmt.Sprintf(`{"txhex":"%s"}`, txHex)
	postData := []byte(stringVal)

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/raw
	txID, err = c.Request("tx/raw", "POST", nil, postData)
	if err != nil {
		return
	}

	// Got an error
	if c.LastRequest.StatusCode > 200 {
		err = fmt.Errorf("error broadcasting: %s", txID)
		txID = "" // remove the error message
		return
	}

	return
}
