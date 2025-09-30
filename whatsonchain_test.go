package whatsonchain

const (
	testKey = "test-key-for-woc-api"
)

// newMockClient returns a client for mocking
func newMockClient(httpClient HTTPInterface) ClientInterface {
	opts := ClientDefaultOptions()
	opts.APIKey = testKey
	return NewClient(NetworkTest, opts, httpClient)
}

// newMockClientBSV returns a BSV client for mocking
func newMockClientBSV(httpClient HTTPInterface) ClientInterface {
	opts := ClientDefaultOptions()
	opts.APIKey = testKey
	return NewClientWithChain(ChainBSV, NetworkTest, opts, httpClient)
}

// newMockClientBTC returns a BTC client for mocking
func newMockClientBTC(httpClient HTTPInterface) ClientInterface {
	opts := ClientDefaultOptions()
	opts.APIKey = testKey
	return NewClientWithChain(ChainBTC, NetworkTest, opts, httpClient)
}
