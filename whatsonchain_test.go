package whatsonchain

import (
	"fmt"
	"log"
	"testing"
)

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client.Parameters.Network != NetworkMain {
		t.Fatal("expected value to be main")
	}
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client, _ := NewClient()
	fmt.Println(client.Parameters.Network)
	// Output:main
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient()
	}
}

// TestClient_GetHealth tests the GetHealth()
func TestClient_GetHealth(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp string
	resp, err = client.GetHealth()
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}
	if resp != "Whats On Chain" {
		t.Fatal("expected value was wrong", resp)
	}

}

// ExampleClient_GetHealth example using GetHealth()
func ExampleClient_GetHealth() {
	client, _ := NewClient()
	resp, _ := client.GetHealth()
	fmt.Println(resp)
	// Output:Whats On Chain
}

// TestClient_GetChainInfo tests the GetChainInfo()
func TestClient_GetChainInfo(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *ChainInfo
	resp, err = client.GetChainInfo()
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp.BestBlockHash) == 0 {
		t.Fatal("failed to get best block hash")
	}

}

// ExampleClient_GetChainInfo example using GetChainInfo()
func ExampleClient_GetChainInfo() {
	client, _ := NewClient()
	resp, _ := client.GetChainInfo()
	log.Println(resp.BestBlockHash)
	fmt.Println("0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9")
	// Output:0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9
}

// TestClient_GetBlockByHash tests the GetBlockByHash()
func TestClient_GetBlockByHash(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockInfo
	hash := "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23"
	resp, err = client.GetBlockByHash(hash)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23" {
		t.Fatal("failed to get the block hash", resp.Hash)
	}

}

// TestClient_GetBlockByHeight tests the GetBlockByHeight()
func TestClient_GetBlockByHeight(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockInfo
	height := int64(604648)
	resp, err = client.GetBlockByHeight(height)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23" {
		t.Fatal("failed to get the block hash", resp.Hash)
	}

}

// TestClient_GetBlockPages tests the GetBlockPages()
func TestClient_GetBlockPages(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockPagesInfo
	hash := "000000000000000000885a4d8e9912f085b42288adc58b3ee5830a7da9f4fef4"
	resp, err = client.GetBlockPages(hash, 1)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(*resp) == 0 {
		t.Fatal("no transactions found", resp)
	}

}

// TestClient_GetTxByHash tests the GetTxByHash()
func TestClient_GetTxByHash(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *TxInfo
	hash := "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96"
	resp, err = client.GetTxByHash(hash)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96" {
		t.Fatal("failed to get the tx hash", resp.Hash)
	}

}

// TestClient_BroadcastTx tests the BroadcastTx()
func TestClient_BroadcastTx(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	txHex := "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000"
	_, err = client.BroadcastTx(txHex)
	if err == nil {
		t.Fatal("error should have occurred")
	}
}

// TestClient_BulkBroadcastTx tests the BulkBroadcastTx()
func TestClient_BulkBroadcastTx(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	txHex1 := "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000"
	txHex2 := "0100000006279fdaff3d61920bc0b017c7b7ed8a8944b1ed1cc1998efc18b71f4828dd3a02010000006a47304402205ba16dccd88461f11e0a705b1015ae84d543acc92cb25012c14578d14cedf677022036bfdd96c6bbe57488b288a1dd3c1c4208154066fb7f8f9dfe59fb6ef53ca01b41210215b9d8176b6697758859e95ff43e61dfea94d44340b16d6f6ac4ae6d61a37ab3ffffffff3efb67c48b0a387f769219ee5c7431eaa49e11a3bfb4e7969d91b0973bcc3a30010000006b483045022100d1a25279f2bb720848717b85c518018c3a26ac584dce2956aeaa4ead86e43b1002200c2e5193fcc2250cc08407b22e7ee73ba48c8755d2aa193b543d01ad97dff85941210306e0678608241a4dc0bde0a39ab3d29dfdef6624ab5a20aa0b821dea85c9d9d4ffffffff9da533cd621baa3d07aca98f07c03115db993bf685afad6a3f321dac61db1731010000006a473044022009769824e84a1b9756aa9450249c19604ffe1474875e9389e851c562dcf272e502206dec4452cdec705bb4066883a2b9cd571768df2017638143bdec4a053b37abcb41210243b052e875c67900cebb6bcf816c7ee510f1a7c1068b0625a1828f1ab82c3806ffffffff4e6647164ee0edc0c7617e82575a88eca16056efb3af297e32fd03767c4e353d010000006a473044022012949fd73deca106804968fd50d86cb89872758f5c492b0938fdb8cee28b30df0220168f88d65400c68c5df348d60c800b1f9fb7e3135af05addac049c940da36bf9412102645b8993f1a183e9d37ec2fd9c5f950110b404aea7e9f01687ab42eb6f8b3563ffffffff49236581eb2ada13ac5972e7211771539ebe06061ac6a27a58acde5521b80949010000006a4730440220621f6dff81bfc0ab6c2b9e321c920f477f4ee5e5a01e3d24aecbb265ad264a1202205e99ac2fe08f6cecfcc8ddc98de09d4693c7f81f03bb478d1ee5936919e00b0d4121039925d3b23e560bd021d665e56f912735fcbb96303051b037a6816a55697a226effffffff1b6196327f18bbe6af3fabd46a5d3af568cc7c4447a340f4949824b6762a98f2010000006b483045022100f679c4408b7d893a6bbfb042685053369de044d4d12c6ddcc678d223684f7a320220669bd27d9e9e2be82c0f653d95cd74a4e188a66b976c5be7ba7d3b113692547041210235157690d4fc237b4aceb9e6f91ae0f0058628df36e4450f173ffae2d7d0d49affffffff0253290000000000001976a914dc57ab8a8365a7263fad7491e9a36601a786772388ac37d60000000000001976a914594d5717bd8f9ae5ae8c56646042082d3d6995f988ac00000000"

	hexSlice := []string{txHex1, txHex2}
	var response *BulkBroadcastResponse
	response, err = client.BulkBroadcastTx(hexSlice, true)
	if err != nil {
		t.Fatal("error  occurred", err)
	}

	if len(response.StatusURL) == 0 {
		t.Fatal("expected a url back", response.StatusURL, response.Feedback)
	}

	if !response.Feedback {
		t.Log("expected to be true")
	}
}

// TestClient_AddressInfo tests the AddressInfo()
func TestClient_AddressInfo(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *AddressInfo
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	resp, err = client.AddressInfo(address)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Address != address {
		t.Fatal("failed to get the address", address, resp.Address)
	}

}

// TestClient_AddressBalance tests the AddressBalance()
func TestClient_AddressBalance(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	//var resp *AddressBalance
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	_, err = client.AddressBalance(address)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	// todo: not so easy to test for unless you leave a wallet with unspent transactions just for this test ;-)
	/*if resp.Unconfirmed != 0 {
		t.Fatal("failed to get the unconfirmed value", resp.Unconfirmed)
	}

	if resp.Confirmed != 0 {
		t.Fatal("failed to get the confirmed value", resp.Unconfirmed)
	}*/

}

// TestClient_AddressHistory tests the AddressHistory()
func TestClient_AddressHistory(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp AddressHistory
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	resp, err = client.AddressHistory(address)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("failed to get history values", resp)
	}

	// todo: this might change! not the best test :P
	if resp[0].TxHash != "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1" {
		t.Fatal("failed to get correct hash", resp, resp[0].TxHash)
	}

}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactions(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	//var resp AddressHistory
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	_, err = client.AddressUnspentTransactions(address)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	// todo: not so easy to test for unless you leave a wallet with unspent transactions just for this test ;-)
	/*if len(resp) == 0 {
		t.Fatal("failed to get history values", resp)
	}*/

}

// TestClient_GetMerkleProof tests the GetMerkleProof()
func TestClient_GetMerkleProof(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *MerkleInfo
	tx := "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96"
	resp, err = client.GetMerkleProof(tx)
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.BlockHeight != 575191 {
		t.Fatal("failed to get the block height", resp.BlockHeight)
	}

	if resp.Merkle[0] != "7e0ba1980522125f1f40d19a249ab3ae036001b991776813d25aebe08e8b8a50" {
		t.Fatal("failed to get the merkle hash", resp.Merkle[0])
	}

}
