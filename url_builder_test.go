package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildURL_Default verifies the default base endpoint is used when no
// WithBaseURL option is provided.
func TestBuildURL_Default(t *testing.T) {
	client, err := NewClient(context.Background())
	require.NoError(t, err)

	c, ok := client.(*Client)
	require.True(t, ok)

	assert.Equal(t, apiEndpointBase, client.BaseURL())
	assert.Equal(t,
		"https://api.whatsonchain.com/v1/bsv/main/tx/abc/hex",
		c.buildURL("/tx/%s/hex", "abc"),
	)
}

// TestWithBaseURL verifies the WithBaseURL option overrides the base endpoint
// and normalizes a missing trailing slash, while an empty value is ignored.
func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		wantBase string
		wantURL  string
	}{
		{
			name:     "custom with trailing slash",
			baseURL:  "https://example.com/v1/",
			wantBase: "https://example.com/v1/",
			wantURL:  "https://example.com/v1/bsv/main/tx/abc/hex",
		},
		{
			name:     "custom without trailing slash is normalized",
			baseURL:  "https://example.com/v1",
			wantBase: "https://example.com/v1/",
			wantURL:  "https://example.com/v1/bsv/main/tx/abc/hex",
		},
		{
			name:     "empty leaves default unchanged",
			baseURL:  "",
			wantBase: apiEndpointBase,
			wantURL:  "https://api.whatsonchain.com/v1/bsv/main/tx/abc/hex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(context.Background(), WithBaseURL(tt.baseURL))
			require.NoError(t, err)

			c, ok := client.(*Client)
			require.True(t, ok)

			assert.Equal(t, tt.wantBase, client.BaseURL())
			assert.Equal(t, tt.wantURL, c.buildURL("/tx/%s/hex", "abc"))
		})
	}
}

// TestWithBaseURL_Network verifies the configured base URL is combined with the
// selected network prefix.
func TestWithBaseURL_Network(t *testing.T) {
	client, err := NewClient(context.Background(),
		WithBaseURL("https://example.com/v1/"),
		WithNetwork(NetworkTest),
	)
	require.NoError(t, err)

	c, ok := client.(*Client)
	require.True(t, ok)

	assert.Equal(t, "https://example.com/v1/bsv/test/block/hash/abc", c.buildURL("/block/hash/%s", "abc"))
}

// TestWithBaseURL_EndToEnd verifies requests are actually sent to the configured
// base URL by pointing the client at a local test server.
func TestWithBaseURL_EndToEnd(t *testing.T) {
	const rawHex = "0100000000"

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = io.WriteString(w, rawHex)
	}))
	defer server.Close()

	client, err := NewClient(context.Background(), WithBaseURL(server.URL+"/v1/"))
	require.NoError(t, err)

	got, err := client.GetRawTransactionData(context.Background(), "abc123")
	require.NoError(t, err)
	assert.Equal(t, rawHex, got)
	assert.Equal(t, "/v1/bsv/main/tx/abc123/hex", gotPath)
}
