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
