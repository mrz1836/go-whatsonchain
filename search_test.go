package whatsonchain

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPSearchValid for mocking requests
type mockHTTPSearchValid struct{}

// queryData is the query for searching
type queryData struct {
	Query string `json:"query"`
}

// Do is a mock http request
func (m *mockHTTPSearchValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	decoder := json.NewDecoder(req.Body)
	var data queryData
	err := decoder.Decode(&data)
	if err != nil {
		return resp, err
	}

	// Valid (address)
	if strings.Contains(data.Query, "1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"results":[{"type":"address","url":"https://whatsonchain.com/address/1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H"}]}`)))
	}

	// Valid (tx)
	if strings.Contains(data.Query, "6a7c821fd13c5cec773f7e221479651804197866469e92a4d6d47e1fd34d090d") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"results":[{"type":"tx","url":"https://whatsonchain.com/tx/6a7c821fd13c5cec773f7e221479651804197866469e92a4d6d47e1fd34d090d"}]}`)))
	}

	// Valid (block)
	if strings.Contains(data.Query, "000000000000000002080d0ad78d08691d956d08fb8556339b6dd84fbbfdf1bc") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"results":[{"type":"block","url":"https://whatsonchain.com/block/000000000000000002080d0ad78d08691d956d08fb8556339b6dd84fbbfdf1bc"}]}`)))
	}

	// Valid (op_return)
	if strings.Contains(data.Query, "unknown") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"results":[{"type":"op_return","url":"https://whatsonchain.com/opreturn-query?term=unknown\u0026size=10\u0026offset=0"}]}`)))
	}

	// Invalid
	if strings.Contains(data.Query, "error") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Not found
	if strings.Contains(data.Query, "notFound") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetExplorerLinks tests the GetExplorerLinks()
func TestClient_GetExplorerLinks(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPSearchValid{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		typeName      string
		url           string
		expectedError bool
		statusCode    int
	}{
		{"1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H", "address", "https://whatsonchain.com/address/1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H", false, http.StatusOK},
		{"6a7c821fd13c5cec773f7e221479651804197866469e92a4d6d47e1fd34d090d", "tx", "https://whatsonchain.com/tx/6a7c821fd13c5cec773f7e221479651804197866469e92a4d6d47e1fd34d090d", false, http.StatusOK},
		{"000000000000000002080d0ad78d08691d956d08fb8556339b6dd84fbbfdf1bc", "block", "https://whatsonchain.com/block/000000000000000002080d0ad78d08691d956d08fb8556339b6dd84fbbfdf1bc", false, http.StatusOK},
		{"unknown", "op_return", "https://whatsonchain.com/opreturn-query?term=unknown&size=10&offset=0", false, http.StatusOK},
		{"error", "", "", true, http.StatusBadRequest},
		{"notFound", "", "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.GetExplorerLinks(ctx, test.input)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && err == nil && output.Results != nil {
			if output.Results[0].Type != test.typeName {
				t.Errorf("%s Failed: [%s] inputted and [%s] type expected, received: [%s]", t.Name(), test.input, test.typeName, output.Results[0].Type)
			}
			if output.Results[0].URL != test.url {
				t.Errorf("%s Failed: [%s] inputted and [%s] url expected, received: [%s]", t.Name(), test.input, test.url, output.Results[0].URL)
			}
		}
	}
}
