package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPScript for mocking requests
type mockHTTPScript struct{}

// Do is a mock http request
func (m *mockHTTPScript) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	// Valid
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"tx_hash":"52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48","height":620539},{"tx_hash":"4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0","height":620539}]`)))
	}

	// Invalid
	if strings.Contains(req.URL.String(), "script/invalidTx/history") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Valid (has utxo)
	if strings.Contains(req.URL.String(), "script/92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"height": 640558,"tx_pos": 1,"tx_hash": "5c6ac3a685be0791aa6e6eedb03d48cbf76046ea499e0a9cefbdc0fb3969ad13","value": 533335}]`)))
	}

	// Valid (empty)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
	}

	// Invalid
	if strings.Contains(req.URL.String(), "script/invalidTx/unspent") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetScriptHistory tests the GetScriptHistory()
func TestClient_GetScriptHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})

	// Create the list of tests
	var tests = []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 620539, "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48", false, http.StatusOK},
		{"invalidTx", 0, "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetScriptHistory(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && output[0].Height != test.height && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
		} else if output != nil && output[0].TxHash != test.hash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_GetScriptUnspentTransactions tests the GetScriptUnspentTransactions()
func TestClient_GetScriptUnspentTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})

	// Create the list of tests
	var tests = []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb", 640558, "5c6ac3a685be0791aa6e6eedb03d48cbf76046ea499e0a9cefbdc0fb3969ad13", false, http.StatusOK},
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 0, "", false, http.StatusOK},
		{"invalidTx", 0, "", true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.GetScriptUnspentTransactions(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && output[0].Height != test.height && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
		} else if output != nil && len(output) > 0 && output[0].TxHash != test.hash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}
