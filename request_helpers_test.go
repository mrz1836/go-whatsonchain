package whatsonchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

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

// errTestNotFound is a not-found sentinel used to test checkStatusCode in isolation
var errTestNotFound = errors.New("test resource not found")

// TestCheckStatusCode covers every branch of checkStatusCode: success, a 404 with
// and without a not-found sentinel, and other failure statuses.
func TestCheckStatusCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		status      int
		body        string
		notFoundErr error
		wantErr     string // exact message; empty means no error
		wantIs      []error
		wantNotIs   []error
		wantBare    bool // the error must be notFoundErr itself, so == comparisons work
	}{
		{
			name:        "200 with body is success",
			status:      http.StatusOK,
			body:        `{"rate":38.5}`,
			notFoundErr: errTestNotFound,
		},
		{
			name:   "200 with empty body is success (emptiness is the caller's concern)",
			status: http.StatusOK,
		},
		{
			name:        "404 with empty body returns the bare sentinel",
			status:      http.StatusNotFound,
			notFoundErr: errTestNotFound,
			wantErr:     "test resource not found",
			wantIs:      []error{errTestNotFound},
			wantNotIs:   []error{ErrRequestFailed},
			wantBare:    true,
		},
		{
			name:        "404 with whitespace-only body returns the bare sentinel",
			status:      http.StatusNotFound,
			body:        " \r\n\t",
			notFoundErr: errTestNotFound,
			wantErr:     "test resource not found",
			wantIs:      []error{errTestNotFound},
			wantNotIs:   []error{ErrRequestFailed},
			wantBare:    true,
		},
		{
			name:        "404 with JSON body wraps the sentinel",
			status:      http.StatusNotFound,
			body:        `{"error":"not found"}`,
			notFoundErr: errTestNotFound,
			wantErr:     `test resource not found: HTTP 404: {"error":"not found"}`,
			wantIs:      []error{errTestNotFound},
			wantNotIs:   []error{ErrRequestFailed},
		},
		{
			name:        "404 for an unknown route wraps the sentinel and keeps the text",
			status:      http.StatusNotFound,
			body:        "404 page not found\n",
			notFoundErr: errTestNotFound,
			wantErr:     "test resource not found: HTTP 404: 404 page not found",
			wantIs:      []error{errTestNotFound},
			wantNotIs:   []error{ErrRequestFailed},
		},
		{
			name:      "404 without a sentinel is a request failure",
			status:    http.StatusNotFound,
			wantErr:   "API request failed: HTTP 404",
			wantIs:    []error{ErrRequestFailed},
			wantNotIs: []error{errTestNotFound},
		},
		{
			name:    "404 without a sentinel keeps the body",
			status:  http.StatusNotFound,
			body:    "404 page not found\n",
			wantErr: "API request failed: HTTP 404: 404 page not found",
			wantIs:  []error{ErrRequestFailed},
		},
		{
			name:        "401 ignores the not-found sentinel",
			status:      http.StatusUnauthorized,
			body:        "Unauthorized",
			notFoundErr: errTestNotFound,
			wantErr:     "API request failed: HTTP 401: Unauthorized",
			wantIs:      []error{ErrRequestFailed},
			wantNotIs:   []error{errTestNotFound},
		},
		{
			name:        "429 with empty body",
			status:      http.StatusTooManyRequests,
			notFoundErr: errTestNotFound,
			wantErr:     "API request failed: HTTP 429",
			wantIs:      []error{ErrRequestFailed},
			wantNotIs:   []error{errTestNotFound},
		},
		{
			name:    "500 trims surrounding whitespace from the body",
			status:  http.StatusInternalServerError,
			body:    "\n  Internal Server Error  \n",
			wantErr: "API request failed: HTTP 500: Internal Server Error",
			wantIs:  []error{ErrRequestFailed},
		},
		{
			name:    "a status of 0 (no response) is a request failure",
			status:  0,
			wantErr: "API request failed: HTTP 0",
			wantIs:  []error{ErrRequestFailed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkStatusCode(tt.status, []byte(tt.body), tt.notFoundErr)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
			for _, target := range tt.wantIs {
				require.ErrorIs(t, err, target)
			}
			for _, target := range tt.wantNotIs {
				require.NotErrorIs(t, err, target)
			}
			if tt.wantBare {
				assert.Same(t, tt.notFoundErr, err, "an empty 404 must return the sentinel itself")
			}
		})
	}
}

// TestEndpoints_HTTP404 checks how every kind of endpoint reports a 404, through
// the public client. The API answers an unknown resource with an empty 404 and an
// unknown route with "404 page not found"; endpoints with a not-found sentinel
// return it (bare for an empty body, wrapped otherwise), and the rest return
// ErrRequestFailed. Before issue #223, string endpoints returned the 404 body as
// data and bulk endpoints returned (nil, nil).
func TestEndpoints_HTTP404(t *testing.T) {
	t.Parallel()

	hashes := &TxHashes{TxIDs: []string{testTxID1}}

	tests := []struct {
		name     string
		call     func(ctx context.Context, c ClientInterface) (any, error)
		notFound error   // nil: the 404 must be reported as ErrRequestFailed
		alsoIs   []error // further sentinels the error must match
	}{
		// String endpoints
		{"GetRawTransactionData", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetRawTransactionData(ctx, testTxID1)
		}, ErrTransactionNotFound, nil},
		{"GetRawTransactionOutputData", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetRawTransactionOutputData(ctx, testTxID1, 0)
		}, ErrTransactionNotFound, nil},
		{"GetOpReturnData", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetOpReturnData(ctx, testTxID1)
		}, ErrTransactionNotFound, nil},
		{"GetTransactionAsBinary", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetTransactionAsBinary(ctx, testTxID1)
		}, ErrTransactionNotFound, nil},
		{"DownloadReceipt", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.DownloadReceipt(ctx, testTxID1)
		}, ErrTransactionNotFound, nil},
		{"DownloadStatement", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.DownloadStatement(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		}, ErrAddressNotFound, nil},
		{"GetLatestHeaderBytes", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetLatestHeaderBytes(ctx, 5)
		}, ErrHeadersNotFound, nil},
		{"GetCirculatingSupply", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetCirculatingSupply(ctx)
		}, ErrChainInfoNotFound, nil},
		{"GetScriptUsed", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetScriptUsed(ctx, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3")
		}, ErrScriptNotFound, nil},
		{"GetHealth", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetHealth(ctx)
		}, nil, nil},

		// JSON endpoints
		{"GetExchangeRate", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetExchangeRate(ctx)
		}, ErrExchangeRateNotFound, nil},
		{"GetTxByHash", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetTxByHash(ctx, testTxID1)
		}, ErrTransactionNotFound, nil},
		{"GetBlockByHeight", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetBlockByHeight(ctx, 999999999)
		}, ErrBlockNotFound, nil},
		{"GetChainTips", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetChainTips(ctx)
		}, ErrChainTipsNotFound, nil},
		{"GetScriptUnspentTransactions", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.GetScriptUnspentTransactions(ctx, "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3")
		}, ErrScriptNotFound, nil},

		// Bulk POST endpoints never use 404 for unknown items, so a 404 is a failure
		{"BulkTransactionDetails", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BulkTransactionDetails(ctx, hashes)
		}, nil, nil},
		{"BulkRawTransactionData", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BulkRawTransactionData(ctx, hashes)
		}, nil, nil},
		{"BulkTransactionStatus", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BulkTransactionStatus(ctx, hashes)
		}, nil, nil},
		{"BulkSpentOutputs", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BulkSpentOutputs(ctx, &BulkSpentOutputRequest{UTXOs: []BulkSpentUTXO{{TxID: testTxID1, Vout: 0}}})
		}, nil, nil},

		// Broadcasts: a 404 means the request never reached a node
		{"BroadcastTx", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BroadcastTx(ctx, "0100000001")
		}, nil, []error{ErrBroadcastFailed}},
		{"BulkBroadcastTx", func(ctx context.Context, c ClientInterface) (any, error) {
			return c.BulkBroadcastTx(ctx, []string{"0100000001"}, true)
		}, nil, []error{ErrBroadcastFailed}},
	}

	bodies := []struct {
		name string
		body string
	}{
		{"empty body (unknown resource)", ""},
		{"text body (unknown route)", "404 page not found\n"},
		{"JSON body", `{"error":"not found"}`},
	}

	for _, tt := range tests {
		for _, b := range bodies {
			t.Run(tt.name+"/"+b.name, func(t *testing.T) {
				t.Parallel()

				client := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusNotFound, body: b.body})
				result, err := tt.call(context.Background(), client)

				require.Error(t, err)
				assert.Empty(t, result, "no data may be returned with a 404")
				assert.Equal(t, http.StatusNotFound, client.LastRequest().StatusCode)
				for _, target := range tt.alsoIs {
					require.ErrorIs(t, err, target)
				}

				if tt.notFound == nil {
					require.ErrorIs(t, err, ErrRequestFailed)
					assert.Contains(t, err.Error(), "HTTP 404")
					return
				}

				require.ErrorIs(t, err, tt.notFound)
				require.NotErrorIs(t, err, ErrRequestFailed)
				if b.body == "" {
					assert.Same(t, tt.notFound, err, "an empty 404 must return the sentinel itself")
				} else {
					assert.Contains(t, err.Error(), "HTTP 404: "+strings.TrimSpace(b.body))
				}
			})
		}
	}
}

// TestRequestHelpers_EmptyAndMalformed200 checks the 200 paths that the 404
// changes must leave alone: an empty body still returns the endpoint's sentinel,
// malformed JSON is still a decode error, and a string endpoint returns an empty
// body as is.
func TestRequestHelpers_EmptyAndMalformed200(t *testing.T) {
	t.Parallel()

	t.Run("empty object response", func(t *testing.T) {
		t.Parallel()

		info, err := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK}).GetChainInfo(context.Background())
		assert.Same(t, ErrChainInfoNotFound, err)
		assert.Nil(t, info)
	})

	t.Run("empty slice response", func(t *testing.T) {
		t.Parallel()

		tips, err := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK}).GetChainTips(context.Background())
		assert.Same(t, ErrChainTipsNotFound, err)
		assert.Nil(t, tips)
	})

	t.Run("malformed object response", func(t *testing.T) {
		t.Parallel()

		info, err := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK, body: "{"}).GetChainInfo(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrChainInfoNotFound)
		require.NotErrorIs(t, err, ErrRequestFailed)
		assert.Nil(t, info)
	})

	t.Run("malformed slice response", func(t *testing.T) {
		t.Parallel()

		tips, err := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK, body: "{"}).GetChainTips(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrChainTipsNotFound)
		require.NotErrorIs(t, err, ErrRequestFailed)
		assert.Nil(t, tips)
	})

	t.Run("empty string response", func(t *testing.T) {
		t.Parallel()

		health, err := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK}).GetHealth(context.Background())
		require.NoError(t, err)
		assert.Empty(t, health)
	})
}

// TestErrorBodySnippet checks that error messages quote at most maxErrorBodyLength
// bytes of a response body, never split a UTF-8 character, and say how much was cut.
func TestErrorBodySnippet(t *testing.T) {
	t.Parallel()

	htmlPage := "<html><body>" + strings.Repeat("<p>Service temporarily unavailable</p>", 1400) + "</body></html>"
	require.Greater(t, len(htmlPage), 50_000)

	tests := []struct {
		name       string
		body       []byte
		wantPrefix string // the snippet must start with this
		wantExact  string // when set, the snippet must equal this
		wantTotal  int    // when set, the body length reported as truncated
	}{
		{name: "nil body", body: nil, wantExact: ""},
		{name: "whitespace only", body: []byte(" \n\t\r\n "), wantExact: ""},
		{name: "short body is trimmed", body: []byte("  Unauthorized \n"), wantExact: "Unauthorized"},
		{name: "exactly the limit is kept whole", body: []byte(strings.Repeat("a", maxErrorBodyLength)), wantExact: strings.Repeat("a", maxErrorBodyLength)},
		{name: "limit after trimming is kept whole", body: []byte("\n" + strings.Repeat("a", maxErrorBodyLength) + "\n"), wantExact: strings.Repeat("a", maxErrorBodyLength)},
		{
			name:      "one byte over the limit is truncated",
			body:      []byte(strings.Repeat("a", maxErrorBodyLength+1)),
			wantExact: strings.Repeat("a", maxErrorBodyLength) + "... (truncated, 513 bytes total)",
			wantTotal: maxErrorBodyLength + 1,
		},
		{
			name:       "large HTML page is truncated",
			body:       []byte(htmlPage),
			wantPrefix: htmlPage[:maxErrorBodyLength],
			wantTotal:  len(htmlPage),
		},
		{
			name:       "total counts the untrimmed body",
			body:       []byte("\n\n" + strings.Repeat("b", 600) + "\n\n"),
			wantPrefix: strings.Repeat("b", maxErrorBodyLength),
			wantTotal:  604,
		},
		{
			// "é" is 2 bytes and straddles the cut: it must be dropped, not split
			name:      "two-byte rune across the cut",
			body:      []byte(strings.Repeat("a", maxErrorBodyLength-1) + "é" + "tail"),
			wantExact: strings.Repeat("a", maxErrorBodyLength-1) + "... (truncated, 517 bytes total)",
			wantTotal: maxErrorBodyLength + 5,
		},
		{
			// "😀" is 4 bytes, starting 2 bytes before the cut
			name:      "four-byte rune across the cut",
			body:      []byte(strings.Repeat("a", maxErrorBodyLength-2) + "😀" + "tail"),
			wantExact: strings.Repeat("a", maxErrorBodyLength-2) + "... (truncated, 518 bytes total)",
			wantTotal: maxErrorBodyLength + 6,
		},
		{
			// "é" ends exactly at the cut, so it is kept
			name:      "rune ending at the cut is kept",
			body:      []byte(strings.Repeat("a", maxErrorBodyLength-2) + "é" + "tail"),
			wantExact: strings.Repeat("a", maxErrorBodyLength-2) + "é... (truncated, 516 bytes total)",
			wantTotal: maxErrorBodyLength + 4,
		},
		{
			// Invalid UTF-8 (all continuation bytes) must not scan back through the whole body
			name:      "binary body is cut close to the limit",
			body:      bytes.Repeat([]byte{0x80}, 2048),
			wantTotal: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := errorBodySnippet(tt.body)
			if tt.wantExact != "" || tt.wantTotal == 0 {
				assert.Equal(t, tt.wantExact, got)
			}
			assert.True(t, strings.HasPrefix(got, tt.wantPrefix), "snippet should start with %q", tt.wantPrefix)

			if tt.wantTotal == 0 {
				assert.NotContains(t, got, "truncated")
				return
			}

			suffix := fmt.Sprintf("... (truncated, %d bytes total)", tt.wantTotal)
			require.True(t, strings.HasSuffix(got, suffix), "snippet %q should end with %q", got, suffix)
			kept := strings.TrimSuffix(got, suffix)
			assert.LessOrEqual(t, len(kept), maxErrorBodyLength)
			assert.GreaterOrEqual(t, len(kept), maxErrorBodyLength-(utf8.UTFMax-1))
			if utf8.Valid(tt.body) {
				assert.True(t, utf8.ValidString(kept), "a valid UTF-8 body must not be split mid-rune")
			}
		})
	}
}

// TestErrorMessages_TruncateLargeBodies checks, through the public client, that an
// error quotes only the start of a large body (issue #223: a 50 MB body used to be
// quoted in full).
func TestErrorMessages_TruncateLargeBodies(t *testing.T) {
	t.Parallel()

	htmlPage := "<html><body>" + strings.Repeat("<p>Bad gateway</p>", 3000) + "<p>END-OF-PAGE</p></body></html>"

	tests := []struct {
		name   string
		status int
		call   func(ctx context.Context, c ClientInterface) error
		is     error
	}{
		{"JSON endpoint 502", http.StatusBadGateway, func(ctx context.Context, c ClientInterface) error {
			_, err := c.GetTxByHash(ctx, testTxID1)
			return err
		}, ErrRequestFailed},
		{"string endpoint 500", http.StatusInternalServerError, func(ctx context.Context, c ClientInterface) error {
			_, err := c.GetRawTransactionData(ctx, testTxID1)
			return err
		}, ErrRequestFailed},
		{"string endpoint 404 with page", http.StatusNotFound, func(ctx context.Context, c ClientInterface) error {
			_, err := c.GetRawTransactionData(ctx, testTxID1)
			return err
		}, ErrTransactionNotFound},
		{"broadcast 400", http.StatusBadRequest, func(ctx context.Context, c ClientInterface) error {
			_, err := c.BroadcastTx(ctx, "0100000001")
			return err
		}, ErrBroadcastFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newMockClient(&mockHTTPStatusCode{statusCode: tt.status, body: htmlPage})
			err := tt.call(context.Background(), client)

			require.ErrorIs(t, err, tt.is)
			msg := err.Error()
			assert.Contains(t, msg, fmt.Sprintf("HTTP %d: <html><body>", tt.status))
			assert.Contains(t, msg, fmt.Sprintf("... (truncated, %d bytes total)", len(htmlPage)))
			assert.NotContains(t, msg, "END-OF-PAGE", "the end of the body must not be quoted")
			assert.Less(t, len(msg), maxErrorBodyLength+200, "the message must stay small")
		})
	}
}

// TestBroadcastTx_ClassifiesFullBody checks that broadcast rejection reasons are
// still recognized when they appear past the truncation point of the error message.
func TestBroadcastTx_ClassifiesFullBody(t *testing.T) {
	t.Parallel()

	body := strings.Repeat("x", 2*maxErrorBodyLength) + " 257: txn-already-known"
	client := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusBadRequest, body: body})

	txID, err := client.BroadcastTx(context.Background(), "0100000001")
	require.ErrorIs(t, err, ErrTxAlreadyInMempool)
	require.ErrorIs(t, err, ErrBroadcastFailed)
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Empty(t, txID)
	assert.NotContains(t, err.Error(), "txn-already-known", "the reason is past the truncation point")
	assert.Contains(t, err.Error(), "truncated")
}
