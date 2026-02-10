package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPStatusCode returns a mock HTTP client that always returns the given status code and body.
type mockHTTPStatusCode struct {
	statusCode int
	body       string
}

func (m *mockHTTPStatusCode) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
	}, nil
}

func TestCheckStatusCode_OK(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusOK,
		body:       `{"confirmed": 100}`,
	})

	// Make a request so LastRequest is populated
	balance, err := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.NoError(t, err)
	require.NotNil(t, balance)
	assert.Equal(t, int64(100), balance.Confirmed)
}

func TestCheckStatusCode_Unauthorized(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusUnauthorized,
		body:       `Unauthorized`,
	})

	balance, err := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.Error(t, err)
	assert.Nil(t, balance)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Contains(t, err.Error(), "401")
	assert.Contains(t, err.Error(), "Unauthorized")
}

func TestCheckStatusCode_RateLimited(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusTooManyRequests,
		body:       `Rate limit exceeded`,
	})

	balance, err := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.Error(t, err)
	assert.Nil(t, balance)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Contains(t, err.Error(), "429")
}

func TestCheckStatusCode_ServerError(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusInternalServerError,
		body:       `Internal Server Error`,
	})

	balance, err := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.Error(t, err)
	assert.Nil(t, balance)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Contains(t, err.Error(), "500")
}

func TestCheckStatusCode_EmptyBody(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusForbidden,
		body:       ``,
	})

	balance, err := client.AddressBalance(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.Error(t, err)
	assert.Nil(t, balance)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Contains(t, err.Error(), "403")
}

func TestCheckStatusCode_SliceEndpoint_Unauthorized(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusUnauthorized,
		body:       `Unauthorized`,
	})

	history, err := client.AddressUnspentTransactions(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.Error(t, err)
	assert.Nil(t, history)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Contains(t, err.Error(), "401")
}

func TestCheckStatusCode_SliceEndpoint_OK(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPStatusCode{
		statusCode: http.StatusOK,
		body:       `{"address":"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","script":"","result":[{"tx_hash": "abc123", "tx_pos": 0, "value": 1000, "height": 100}],"error":""}`,
	})

	history, err := client.AddressUnspentTransactions(context.Background(), "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
	require.NoError(t, err)
	require.Len(t, history, 1)
	assert.Equal(t, "abc123", history[0].TxHash)
}
