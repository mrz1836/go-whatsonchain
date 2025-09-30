package whatsonchain

import "context"

const (
	testKey = "test-key-for-woc-api"
)

// newMockClient returns a client for mocking
func newMockClient(httpClient HTTPInterface) ClientInterface {
	client, _ := NewClient(
		context.Background(),
		WithNetwork(NetworkTest),
		WithAPIKey(testKey),
		WithHTTPClient(httpClient),
	)
	return client
}

// newMockClientBSV returns a BSV client for mocking
func newMockClientBSV(httpClient HTTPInterface) ClientInterface {
	client, _ := NewClient(
		context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkTest),
		WithAPIKey(testKey),
		WithHTTPClient(httpClient),
	)
	return client
}

// newMockClientBTC returns a BTC client for mocking
func newMockClientBTC(httpClient HTTPInterface) ClientInterface {
	client, _ := NewClient(
		context.Background(),
		WithChain(ChainBTC),
		WithNetwork(NetworkTest),
		WithAPIKey(testKey),
		WithHTTPClient(httpClient),
	)
	return client
}
