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

	// Valid (chain tips)
	if strings.Contains(req.URL.String(), "/chain/tips") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"height":905064,"hash":"0000000000000000112ea3732c0417a2cee0130e9217dbba1b0ff078c92c904e","branchlen":0,"status":"active"},{"height":905048,"hash":"0000000000000000089a618e0f8b8fa5dcc8a201211d09aea1434a1e03c357cf","branchlen":1,"status":"valid-headers"}]`))
	}

	// Valid (peer info)
	if strings.Contains(req.URL.String(), "/peer/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"id":4,"addr":"99.127.49.102:49294","addrlocal":"135.181.137.155:8333","services":"0000000000000021","relaytxes":true,"lastsend":1759193615,"lastrecv":1759193615,"bytessent":806424356,"bytesrecv":1220025525,"conntime":1758556738,"timeoffset":-35,"pingtime":0.147154,"minping":0.136913,"version":70016,"subver":"/Bitcoin SV:1.1.0/","inbound":true,"addnode":false,"startingheight":915509,"txninvsize":0,"banscore":0,"synced_headers":916565,"synced_blocks":916565,"whitelisted":false,"bytessent_per_msg":1234.56,"bytesrecv_per_msg":2345.67}]`))
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

	// Invalid (chain tips)
	if strings.Contains(req.URL.String(), "/chain/tips") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Invalid (peer info)
	if strings.Contains(req.URL.String(), "/peer/info") {
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

	// Not found (chain tips)
	if strings.Contains(req.URL.String(), "/chain/tips") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Not found (peer info)
	if strings.Contains(req.URL.String(), "/peer/info") {
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

// TestClient_GetChainTips tests the GetChainTips()
func TestClient_GetChainTips(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPChainValid{})
	ctx := context.Background()

	// Test the valid response
	tips, err := client.GetChainTips(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if tips == nil {
		t.Errorf("%s Failed: tips was nil", t.Name())
	} else if len(tips) != 2 {
		t.Errorf("%s Failed: tips length was [%d] expected [%d]", t.Name(), len(tips), 2)
	} else if tips[0].Height != 905064 {
		t.Errorf("%s Failed: first tip height was [%d] expected [%d]", t.Name(), tips[0].Height, 905064)
	} else if tips[0].Status != "active" {
		t.Errorf("%s Failed: first tip status was [%s] expected [%s]", t.Name(), tips[0].Status, "active")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPChainInvalid{})

	// Test response
	_, err = client.GetChainTips(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPChainNotFound{})

	// Test response
	_, err = client.GetChainTips(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_GetPeerInfo tests the GetPeerInfo()
func TestClient_GetPeerInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPChainValid{})
	ctx := context.Background()

	// Test the valid response
	peers, err := client.GetPeerInfo(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if peers == nil {
		t.Errorf("%s Failed: peers was nil", t.Name())
	} else if len(peers) != 1 {
		t.Errorf("%s Failed: peers length was [%d] expected [%d]", t.Name(), len(peers), 1)
	} else if peers[0].ID != 4 {
		t.Errorf("%s Failed: peer ID was [%d] expected [%d]", t.Name(), peers[0].ID, 4)
	} else if peers[0].Version != 70016 {
		t.Errorf("%s Failed: peer version was [%d] expected [%d]", t.Name(), peers[0].Version, 70016)
	} else if peers[0].SubVer != "/Bitcoin SV:1.1.0/" {
		t.Errorf("%s Failed: peer subver was [%s] expected [%s]", t.Name(), peers[0].SubVer, "/Bitcoin SV:1.1.0/")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPChainInvalid{})

	// Test response
	_, err = client.GetPeerInfo(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPChainNotFound{})

	// Test response
	_, err = client.GetPeerInfo(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
