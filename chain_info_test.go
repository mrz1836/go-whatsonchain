package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		return resp, fmt.Errorf("missing request")
	}

	// Valid (chain info)
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"chain":"main","blocks":640504,"headers":640504,"bestblockhash":"00000000000000000187b269ba0ed06be21c0d0d623c68957ad0308b3004f8ee","difficulty":286794300954.8341,"mediantime":1592843022,"verificationprogress":0.9999928741979456,"pruned":false,"chainwork":"0000000000000000000000000000000000000000010e6322afd01e2bb1415909"}`)))
	}

	// Valid (circulating supply)
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`18440650`)))
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
		return resp, fmt.Errorf("missing request")
	}

	// Invalid (chain info)
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Invalid (circulating supply)
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetChainInfo tests the GetChainInfo()
func TestClient_GetChainInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(t, &mockHTTPChainValid{})

	// Test the valid response
	info, err := client.GetChainInfo()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info.Blocks != 640504 {
		t.Errorf("%s Failed: blocks was [%d] expected [%d]", t.Name(), info.Blocks, 640504)
	}

	// New invalid mock client
	client = newMockClient(t, &mockHTTPChainInvalid{})

	// Test invalid response
	_, err = client.GetChainInfo()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// ExampleClient_GetChainInfo example using GetChainInfo()
func ExampleClient_GetChainInfo() {
	client, _ := NewClient(NetworkMain, nil)
	resp, _ := client.GetChainInfo()
	log.Println(resp.BestBlockHash)
	fmt.Println("0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9")
	// Output:0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9
}

// TestClient_GetCirculatingSupply tests the GetCirculatingSupply()
func TestClient_GetCirculatingSupply(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(t, &mockHTTPChainValid{})

	// Test the valid response
	supply, err := client.GetCirculatingSupply()
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if supply != 18440650 {
		t.Errorf("%s Failed: supply was [%f] expected [%d]", t.Name(), supply, 18440650)
	}

	// New invalid mock client
	client = newMockClient(t, &mockHTTPChainInvalid{})

	// Test invalid response
	_, err = client.GetCirculatingSupply()
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// ExampleClient_GetCirculatingSupply example using GetCirculatingSupply()
func ExampleClient_GetCirculatingSupply() {
	client, _ := NewClient(NetworkMain, nil)
	supply, _ := client.GetCirculatingSupply()
	log.Printf("%f", supply)
	fmt.Println("18225787.5")
	// Output:18225787.5
}
