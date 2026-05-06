package whatsonchain

import "context"

const (
	testKey          = "test-key-for-woc-api"
	testAddress1     = "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	testAddress2     = "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP"
	testMockError    = "error"
	testMockNotFound = "notFound"
	testScriptHash1  = "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3"
	testTxID1        = "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96"
	testTxID2        = "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258"
	testTxIDInvalid  = "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1VV"
	testTxID2Invalid = "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba32VV"
	testAddress3     = "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"
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
