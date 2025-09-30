package whatsonchain

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test error variables
var (
	errScriptMissingRequest  = errors.New("missing request")
	errScriptBadRequest      = errors.New("bad request")
	errScriptNoValidResponse = errors.New("no valid response found")
)

// mockHTTPScript for mocking requests
type mockHTTPScript struct{}

// Do is a mock http request
func (m *mockHTTPScript) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errScriptMissingRequest
	}

	// Valid
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"tx_hash":"52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48","height":620539},{"tx_hash":"4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0","height":620539}]`))
	}

	// Invalid (any endpoint with invalidTx)
	if strings.Contains(req.URL.String(), "script/invalidTx/") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, errScriptBadRequest
	}

	// Not found (any endpoint with notFound)
	if strings.Contains(req.URL.String(), "script/notFound/") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Valid (has utxo)
	if strings.Contains(req.URL.String(), "script/92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"height": 640558,"tx_pos": 1,"tx_hash": "5c6ac3a685be0791aa6e6eedb03d48cbf76046ea499e0a9cefbdc0fb3969ad13","value": 533335}]`))
	}

	// Valid (empty)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
	}

	// Invalid
	if strings.Contains(req.URL.String(), "script/invalidTx/unspent/all") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, errScriptBadRequest
	}

	// Not found
	if strings.Contains(req.URL.String(), "script/notFound/unspent/all") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Valid (unspent/all)
	if strings.Contains(req.URL.String(), "/scripts/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"script":"f814a7c3a40164aacc440871e8b7b14eb6a45f0ca7dcbeaea709edc83274c5e7","unspent":[{"height":620539,"tx_pos":0,"tx_hash":"4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0","value":450000}],"error":""},{"script":"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3","unspent":[],"error":""}]`))
	}

	// Valid (script used - true)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/used") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString("true"))
	}

	// Valid (script used - false)
	if strings.Contains(req.URL.String(), "script/92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb/used") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString("false"))
	}

	// Valid (unconfirmed history)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/unconfirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"tx_hash":"52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48","height":620539}]`))
	}

	// Valid (confirmed history)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/confirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"tx_hash":"4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0","height":620539}]`))
	}

	// Valid (bulk unconfirmed history)
	if strings.Contains(req.URL.String(), "/scripts/unconfirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"script":"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3","history":[{"tx_hash":"52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48","height":620539}],"error":""}]`))
	}

	// Valid (bulk confirmed history)
	if strings.Contains(req.URL.String(), "/scripts/confirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"script":"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3","history":[{"tx_hash":"4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0","height":620539}],"error":""}]`))
	}

	// Valid (unconfirmed unspent)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/unconfirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"height": 0,"tx_pos": 1,"tx_hash": "abc123","value": 100000}]`))
	}

	// Empty (unconfirmed unspent)
	if strings.Contains(req.URL.String(), "script/92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb/unconfirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
	}

	// Valid (confirmed unspent)
	if strings.Contains(req.URL.String(), "script/995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3/confirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"height": 620539,"tx_pos": 0,"tx_hash": "def456","value": 200000}]`))
	}

	// Empty (confirmed unspent)
	if strings.Contains(req.URL.String(), "script/92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb/confirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
	}

	// Valid (bulk unconfirmed unspent)
	if strings.Contains(req.URL.String(), "/scripts/unconfirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"script":"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3","unspent":[{"height":0,"tx_pos":1,"tx_hash":"abc123","value":100000}],"error":""}]`))
	}

	// Valid (bulk confirmed unspent)
	if strings.Contains(req.URL.String(), "/scripts/confirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"script":"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3","unspent":[{"height":620539,"tx_pos":0,"tx_hash":"def456","value":200000}],"error":""}]`))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPScriptErrors for mocking requests
type mockHTTPScriptErrors struct{}

// Do is a mock http request
func (m *mockHTTPScriptErrors) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errScriptMissingRequest
	}

	// Invalid (bulk endpoints) return an error
	if strings.Contains(req.URL.String(), "/scripts/unspent/all") ||
		strings.Contains(req.URL.String(), "/scripts/unconfirmed/history") ||
		strings.Contains(req.URL.String(), "/scripts/confirmed/history") ||
		strings.Contains(req.URL.String(), "/scripts/unconfirmed/unspent") ||
		strings.Contains(req.URL.String(), "/scripts/confirmed/unspent") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, errScriptMissingRequest
	}

	return nil, errScriptNoValidResponse
}

// mockHTTPScriptNotFound for mocking requests
type mockHTTPScriptNotFound struct{}

// Do is a mock http request
func (m *mockHTTPScriptNotFound) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusNotFound

	// No req found
	if req == nil {
		return resp, errScriptMissingRequest
	}

	// Always return empty body for not found
	resp.Body = io.NopCloser(bytes.NewBufferString(""))
	return resp, nil
}

// TestClient_GetScriptHistory tests the GetScriptHistory()
func TestClient_GetScriptHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 620539, "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48", false, http.StatusOK},
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.GetScriptHistory(ctx, test.input)

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

		if !test.expectedError && output != nil {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_GetScriptUnspentTransactions tests the GetScriptUnspentTransactions()
func TestClient_GetScriptUnspentTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb", 640558, "5c6ac3a685be0791aa6e6eedb03d48cbf76046ea499e0a9cefbdc0fb3969ad13", false, http.StatusOK},
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 0, "", true, http.StatusOK}, // Empty response should error
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.GetScriptUnspentTransactions(ctx, test.input)

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

		if !test.expectedError && len(output) > 0 {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_BulkScriptUnspentTransactions tests the BulkScriptUnspentTransactions()
func TestClient_BulkScriptUnspentTransactions(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		balances, err := client.BulkScriptUnspentTransactions(ctx, &ScriptsList{Scripts: []string{
			"f814a7c3a40164aacc440871e8b7b14eb6a45f0ca7dcbeaea709edc83274c5e7",
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
	})

	t.Run("max scripts (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		balances, err := client.BulkScriptUnspentTransactions(ctx, &ScriptsList{Scripts: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"16",
			"17",
			"18",
			"19",
			"20",
			"21",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptErrors{})
		ctx := context.Background()
		balances, err := client.BulkScriptUnspentTransactions(ctx, &ScriptsList{Scripts: []string{
			"f814a7c3a40164aacc440871e8b7b14eb6a45f0ca7dcbeaea709edc83274c5e7",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("not found (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptNotFound{})
		ctx := context.Background()
		balances, err := client.BulkScriptUnspentTransactions(ctx, &ScriptsList{Scripts: []string{
			"notFound",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_GetScriptUsed tests the GetScriptUsed()
func TestClient_GetScriptUsed(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		expectedUsed  bool
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", true, false, http.StatusOK},
		{"92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb", false, false, http.StatusOK},
		{"invalidTx", false, true, http.StatusBadRequest},
		{"notFound", false, true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		used, err := client.GetScriptUsed(ctx, test.input)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, used, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && used != test.expectedUsed {
			t.Errorf("%s Failed: [%s] inputted and [%t] expected, received: [%t]", t.Name(), test.input, test.expectedUsed, used)
		}
	}
}

// TestClient_GetScriptUnconfirmedHistory tests the GetScriptUnconfirmedHistory()
func TestClient_GetScriptUnconfirmedHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 620539, "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48", false, http.StatusOK},
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.GetScriptUnconfirmedHistory(ctx, test.input)

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

		if !test.expectedError && output != nil && len(output) > 0 {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_GetScriptConfirmedHistory tests the GetScriptConfirmedHistory()
func TestClient_GetScriptConfirmedHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 620539, "4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0", false, http.StatusOK},
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.GetScriptConfirmedHistory(ctx, test.input)

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

		if !test.expectedError && output != nil && len(output) > 0 {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_BulkScriptUnconfirmedHistory tests the BulkScriptUnconfirmedHistory()
func TestClient_BulkScriptUnconfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		histories, err := client.BulkScriptUnconfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.NoError(t, err)
		assert.NotNil(t, histories)
		assert.Len(t, histories, 1)
		assert.Equal(t, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", histories[0].Script)
		assert.Len(t, histories[0].History, 1)
		assert.Equal(t, "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48", histories[0].History[0].TxHash)
	})

	t.Run("max scripts (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		histories, err := client.BulkScriptUnconfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, histories)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptErrors{})
		ctx := context.Background()
		histories, err := client.BulkScriptUnconfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.Error(t, err)
		assert.Nil(t, histories)
	})
}

// TestClient_BulkScriptConfirmedHistory tests the BulkScriptConfirmedHistory()
func TestClient_BulkScriptConfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		histories, err := client.BulkScriptConfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.NoError(t, err)
		assert.NotNil(t, histories)
		assert.Len(t, histories, 1)
		assert.Equal(t, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", histories[0].Script)
		assert.Len(t, histories[0].History, 1)
		assert.Equal(t, "4ec3b63d764558303eda720e8e51f69bbcfe81376075657313fb587306f8a9b0", histories[0].History[0].TxHash)
	})

	t.Run("max scripts (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		histories, err := client.BulkScriptConfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, histories)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptErrors{})
		ctx := context.Background()
		histories, err := client.BulkScriptConfirmedHistory(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.Error(t, err)
		assert.Nil(t, histories)
	})
}

// TestClient_ScriptUnconfirmedUTXOs tests the ScriptUnconfirmedUTXOs()
func TestClient_ScriptUnconfirmedUTXOs(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 0, "abc123", false, http.StatusOK},
		{"92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb", 0, "", true, http.StatusOK}, // Empty response should error
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.ScriptUnconfirmedUTXOs(ctx, test.input)

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

		if !test.expectedError && len(output) > 0 {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_BulkScriptUnconfirmedUTXOs tests the BulkScriptUnconfirmedUTXOs()
func TestClient_BulkScriptUnconfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptUnconfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.NoError(t, err)
		assert.NotNil(t, unspentList)
		assert.Len(t, unspentList, 1)
		assert.Equal(t, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", unspentList[0].Script)
		assert.Len(t, unspentList[0].Utxos, 1)
		assert.Equal(t, "abc123", unspentList[0].Utxos[0].TxHash)
		assert.Equal(t, int64(100000), unspentList[0].Utxos[0].Value)
	})

	t.Run("max scripts (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptUnconfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, unspentList)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptErrors{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptUnconfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.Error(t, err)
		assert.Nil(t, unspentList)
	})
}

// TestClient_ScriptConfirmedUTXOs tests the ScriptConfirmedUTXOs()
func TestClient_ScriptConfirmedUTXOs(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPScript{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		height        int64
		hash          string
		expectedError bool
		statusCode    int
	}{
		{"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", 620539, "def456", false, http.StatusOK},
		{"92cf18576a49ddad3e18f4af23b85d8d8218e03ce3b7533aced3fdd286f7e6cb", 0, "", true, http.StatusOK}, // Empty response should error
		{"invalidTx", 0, "", true, http.StatusBadRequest},
		{"notFound", 0, "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.ScriptConfirmedUTXOs(ctx, test.input)

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

		if !test.expectedError && len(output) > 0 {
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].TxHash != test.hash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.hash, output[0].TxHash)
			}
		}
	}
}

// TestClient_BulkScriptConfirmedUTXOs tests the BulkScriptConfirmedUTXOs()
func TestClient_BulkScriptConfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptConfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.NoError(t, err)
		assert.NotNil(t, unspentList)
		assert.Len(t, unspentList, 1)
		assert.Equal(t, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3", unspentList[0].Script)
		assert.Len(t, unspentList[0].Utxos, 1)
		assert.Equal(t, "def456", unspentList[0].Utxos[0].TxHash)
		assert.Equal(t, int64(200000), unspentList[0].Utxos[0].Value)
	})

	t.Run("max scripts (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScript{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptConfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, unspentList)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPScriptErrors{})
		ctx := context.Background()
		unspentList, err := client.BulkScriptConfirmedUTXOs(ctx, &ScriptsList{Scripts: []string{
			"995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3",
		}})
		require.Error(t, err)
		assert.Nil(t, unspentList)
	})
}

// TestClient_GetScriptUsed_EmptyResponse tests GetScriptUsed with empty response
func TestClient_GetScriptUsed_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPScriptNotFound{})
	ctx := context.Background()
	used, err := client.GetScriptUsed(ctx, "notFound")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrScriptNotFound)
	assert.False(t, used)
}
