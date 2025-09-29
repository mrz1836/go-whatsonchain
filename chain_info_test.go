package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPChainValid for mocking requests
type mockHTTPChainValid struct{}

// Do is a mock http request
func (m *mockHTTPChainValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid (chain info)
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"chain":"main","blocks":640504,"headers":640504,"bestblockhash":"00000000000000000187b269ba0ed06be21c0d0d623c68957ad0308b3004f8ee","difficulty":286794300954.8341,"mediantime":1592843022,"verificationprogress":0.9999928741979456,"pruned":false,"chainwork":"0000000000000000000000000000000000000000010e6322afd01e2bb1415909"}`))
	}

	// Valid (circulating supply)
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`18440650`))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPChainInvalid for mocking requests
type mockHTTPChainInvalid struct{}

// Do is a mock http request
func (m *mockHTTPChainInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Invalid (chain info)
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Invalid (circulating supply)
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Default is valid
	return resp, nil
}

// mockHTTPChainNotFound for mocking requests
type mockHTTPChainNotFound struct{}

// Do is a mock http request
func (m *mockHTTPChainNotFound) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusNotFound

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Not found (chain info)
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Not found (circulating supply)
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetChainInfo tests the GetChainInfo()
func TestClient_GetChainInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPChainValid{})
	ctx := context.Background()

	// Test the valid response
	info, err := client.GetChainInfo(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info.Blocks != 640504 {
		t.Errorf("%s Failed: blocks was [%d] expected [%d]", t.Name(), info.Blocks, 640504)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPChainInvalid{})

	// Test response
	_, err = client.GetChainInfo(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPChainNotFound{})

	// Test response
	_, err = client.GetChainInfo(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_GetCirculatingSupply tests the GetCirculatingSupply()
func TestClient_GetCirculatingSupply(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPChainValid{})
	ctx := context.Background()

	// Test the valid response
	supply, err := client.GetCirculatingSupply(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if supply != 18440650 {
		t.Errorf("%s Failed: supply was [%f] expected [%d]", t.Name(), supply, 18440650)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPChainInvalid{})

	// Test response
	_, err = client.GetCirculatingSupply(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPChainNotFound{})

	// Test response
	_, err = client.GetCirculatingSupply(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
