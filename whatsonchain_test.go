package whatsonchain

// newMockClient returns a client for mocking
func newMockClient(httpClient HTTPInterface) *Client {
	client := NewClient(NetworkTest, nil, httpClient)
	// client.httpClient = httpClient
	return client
}
