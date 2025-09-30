package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPOpReturnValid for mocking valid OP_RETURN requests
type mockHTTPOpReturnValid struct{}

// Do is a mock http request
func (m *mockHTTPOpReturnValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid OP_RETURN endpoint
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/opreturn") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`48656c6c6f20576f726c64`))
	}

	return resp, nil
}

// mockHTTPOpReturnNotFound for mocking not found responses
type mockHTTPOpReturnNotFound struct{}

// Do is a mock http request
func (m *mockHTTPOpReturnNotFound) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusNotFound

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	resp.Body = io.NopCloser(bytes.NewBufferString(`{"error":"Transaction not found"}`))
	return resp, ErrTransactionNotFound
}

// mockHTTPOpReturnEmpty for mocking empty responses
type mockHTTPOpReturnEmpty struct{}

// Do is a mock http request
func (m *mockHTTPOpReturnEmpty) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	resp.Body = io.NopCloser(bytes.NewBufferString(``))
	return resp, nil
}

// TestClient_GetOpReturnData tests the GetOpReturnData method
func TestClient_GetOpReturnData(t *testing.T) {
	t.Parallel()

	t.Run("successful OP_RETURN data retrieval", func(t *testing.T) {
		client := newMockClientBSV(&mockHTTPOpReturnValid{})

		data, err := client.GetOpReturnData(context.Background(), "46c5495468b68248b69e55aa76a6b9ca1cb343bee9477c9c121358380e421ff3")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data != "48656c6c6f20576f726c64" {
			t.Errorf("expected '48656c6c6f20576f726c64', got '%s'", data)
		}
	})

	t.Run("error response handling", func(t *testing.T) {
		client := newMockClientBSV(&mockHTTPOpReturnNotFound{})

		_, err := client.GetOpReturnData(context.Background(), "nonexistent")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("empty response handling", func(t *testing.T) {
		client := newMockClientBSV(&mockHTTPOpReturnEmpty{})

		data, err := client.GetOpReturnData(context.Background(), "empty")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data != "" {
			t.Errorf("expected empty string, got '%s'", data)
		}
	})

	t.Run("chain restriction - BTC client", func(t *testing.T) {
		btcClient := newMockClientBTC(&mockHTTPOpReturnValid{})

		_, err := btcClient.GetOpReturnData(context.Background(), "test")
		if err == nil {
			t.Fatal("expected an error for BTC chain, got nil")
		}

		expectedError := "operation is only available for BSV chain"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("chain restriction - BSV client", func(t *testing.T) {
		bsvClient := newMockClientBSV(&mockHTTPOpReturnValid{})

		data, err := bsvClient.GetOpReturnData(context.Background(), "test")
		if err != nil {
			t.Fatalf("BSV client should allow GetOpReturnData, got error: %v", err)
		}

		if data != "" {
			// Should get the mock response
			t.Logf("got data: %s", data)
		}
	})
}
