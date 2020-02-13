package whatsonchain

import "testing"

// TestClient_GetTxByHash tests the GetTxByHash()
func TestClient_GetTxByHash(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *TxInfo
	hash := "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96"
	if resp, err = client.GetTxByHash(hash); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96" {
		t.Fatal("failed to get the tx hash", resp.Hash)
	}

}

// TestClient_GetTxsByHashes tests the GetTxsByHashes()
func TestClient_GetTxsByHashes(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	hashes := new(TxHashes)
	hashes.TxIDs = append(hashes.TxIDs, "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa", "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258")

	var resp TxList
	if resp, err = client.GetTxsByHashes(hashes); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("failed to get transactions")
	}

	if resp[0].TxID != "294cd1ebd5689fdee03509f92c32184c0f52f037d4046af250229b97e0c8f1aa" {
		t.Fatal("failed to get transaction by hash", resp[0].TxID)
	}

	if resp[1].TxID != "91f68c2c598bc73812dd32d60ab67005eac498bef5f0c45b822b3c9468ba3258" {
		t.Fatal("failed to get transaction by hash", resp[1].TxID)
	}

}

// TestClient_BroadcastTx tests the BroadcastTx()
func TestClient_BroadcastTx(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	txHex := "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000"
	if _, err = client.BroadcastTx(txHex); err == nil {
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
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	txHex1 := "0100000001d1bda0bde67183817b21af863adaa31fda8cafcf2083ca1eaba3054496cbde10010000006a47304402205fddd6abab6b8e94f36bfec51ba2e1f3a91b5327efa88264b5530d0c86538723022010e51693e3d52347d4d2ff142b85b460d3953e625d1e062a5fa2569623fb0ea94121029df3723daceb1fef64fa0558371bc48cc3a7a8e35d8e05b87137dc129a9d4598ffffffff0115d40000000000001976a91459cc95a8cde59ceda718dbf70e612dba4034552688ac00000000"
	txHex2 := "0100000006279fdaff3d61920bc0b017c7b7ed8a8944b1ed1cc1998efc18b71f4828dd3a02010000006a47304402205ba16dccd88461f11e0a705b1015ae84d543acc92cb25012c14578d14cedf677022036bfdd96c6bbe57488b288a1dd3c1c4208154066fb7f8f9dfe59fb6ef53ca01b41210215b9d8176b6697758859e95ff43e61dfea94d44340b16d6f6ac4ae6d61a37ab3ffffffff3efb67c48b0a387f769219ee5c7431eaa49e11a3bfb4e7969d91b0973bcc3a30010000006b483045022100d1a25279f2bb720848717b85c518018c3a26ac584dce2956aeaa4ead86e43b1002200c2e5193fcc2250cc08407b22e7ee73ba48c8755d2aa193b543d01ad97dff85941210306e0678608241a4dc0bde0a39ab3d29dfdef6624ab5a20aa0b821dea85c9d9d4ffffffff9da533cd621baa3d07aca98f07c03115db993bf685afad6a3f321dac61db1731010000006a473044022009769824e84a1b9756aa9450249c19604ffe1474875e9389e851c562dcf272e502206dec4452cdec705bb4066883a2b9cd571768df2017638143bdec4a053b37abcb41210243b052e875c67900cebb6bcf816c7ee510f1a7c1068b0625a1828f1ab82c3806ffffffff4e6647164ee0edc0c7617e82575a88eca16056efb3af297e32fd03767c4e353d010000006a473044022012949fd73deca106804968fd50d86cb89872758f5c492b0938fdb8cee28b30df0220168f88d65400c68c5df348d60c800b1f9fb7e3135af05addac049c940da36bf9412102645b8993f1a183e9d37ec2fd9c5f950110b404aea7e9f01687ab42eb6f8b3563ffffffff49236581eb2ada13ac5972e7211771539ebe06061ac6a27a58acde5521b80949010000006a4730440220621f6dff81bfc0ab6c2b9e321c920f477f4ee5e5a01e3d24aecbb265ad264a1202205e99ac2fe08f6cecfcc8ddc98de09d4693c7f81f03bb478d1ee5936919e00b0d4121039925d3b23e560bd021d665e56f912735fcbb96303051b037a6816a55697a226effffffff1b6196327f18bbe6af3fabd46a5d3af568cc7c4447a340f4949824b6762a98f2010000006b483045022100f679c4408b7d893a6bbfb042685053369de044d4d12c6ddcc678d223684f7a320220669bd27d9e9e2be82c0f653d95cd74a4e188a66b976c5be7ba7d3b113692547041210235157690d4fc237b4aceb9e6f91ae0f0058628df36e4450f173ffae2d7d0d49affffffff0253290000000000001976a914dc57ab8a8365a7263fad7491e9a36601a786772388ac37d60000000000001976a914594d5717bd8f9ae5ae8c56646042082d3d6995f988ac00000000"

	hexSlice := []string{txHex1, txHex2}
	var response *BulkBroadcastResponse
	if response, err = client.BulkBroadcastTx(hexSlice, true); err != nil {
		t.Fatal("error  occurred", err)
	}

	if len(response.StatusURL) == 0 {
		t.Fatal("expected a url back", response.StatusURL, response.Feedback)
	}

	if !response.Feedback {
		t.Log("expected to be true")
	}
}

// TestClient_GetMerkleProof tests the GetMerkleProof()
func TestClient_GetMerkleProof(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp MerkleResults
	tx := "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96"
	if resp, err = client.GetMerkleProof(tx); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp == nil {
		t.Fatal("missing results from proof")
	}

	if resp[0].BlockHash != "0000000000000000091216c46973d82db057a6f9911352892b7769ed517681c3" {
		t.Fatal("failed to match block hash", resp[0].BlockHash)
	}

	if resp[0].Hash != "c1d32f28baa27a376ba977f6a8de6ce0a87041157cef0274b20bfda2b0d8df96" {
		t.Fatal("failed to match hash", resp[0].Hash)
	}

	if resp[0].MerkleRoot != "95a920b1002bed05379a0d2650bb13eb216138f28ee80172f4cf21048528dc60" {
		t.Fatal("failed to match merkle root", resp[0].MerkleRoot)
	}

	if len(resp[0].Branches) != 2 {
		t.Fatal("failed to find expected branches", resp[0].Branches)
	}
}
