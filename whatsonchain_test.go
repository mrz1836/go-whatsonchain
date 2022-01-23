package whatsonchain

// newMockClient returns a client for mocking
func newMockClient(httpClient HTTPInterface) ClientInterface {
	return NewClient(NetworkTest, nil, httpClient)
}
