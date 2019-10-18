package whatsonchain

// ChainInfo is the structure response from getting info about the chain
type ChainInfo struct {
	BestBlockHash        string  `json:"bestblockhash"`
	Blocks               int64   `json:"blocks"`
	Chain                string  `json:"chain"`
	ChainWork            string  `json:"chainwork"`
	Difficulty           float64 `json:"difficulty"`
	Headers              int64   `json:"headers"`
	MedianTime           int64   `json:"mediantime"`
	Pruned               bool    `json:"pruned"`
	VerificationProgress float64 `json:"verificationprogress"`
}

// BlockInfo is the response info about a returned block
type BlockInfo struct {
	Bits              string         `json:"bits"`
	ChainWork         string         `json:"chainwork"`
	CoinbaseTx        CoinbaseTxInfo `json:"coinbaseTx"`
	Confirmations     int64          `json:"confirmations"`
	Difficulty        float64        `json:"difficulty"`
	Hash              string         `json:"hash"`
	Height            int64          `json:"height"`
	MedianTime        int64          `json:"mediantime"`
	MerkleRoot        string         `json:"merkleroot"`
	Miner             string         `json:"Bmgpool"`
	NextBlockHash     string         `json:"nextblockhash"`
	Nonce             int64          `json:"nonce"`
	Pages             Page           `json:"pages"`
	PreviousBlockHash string         `json:"previousblockhash"`
	Size              int64          `json:"size"`
	Time              int64          `json:"time"`
	TotalFees         float64        `json:"totalFees"`
	Tx                []string       `json:"tx"`
	TxCount           int64          `json:"txcount"`
	Version           int64          `json:"version"`
	VersionHex        string         `json:"versionHex"`
}

// TxInfo is the response info about a returned tx
type TxInfo struct {
	BlockHash     string     `json:"blockhash"`
	BlockTime     int64      `json:"blocktime"`
	Confirmations int64      `json:"confirmations"`
	Hash          string     `json:"hash"`
	Hex           string     `json:"hex"`
	LockTime      int64      `json:"locktime"`
	Size          int64      `json:"size"`
	Time          int64      `json:"time"`
	TxID          string     `json:"txid"`
	Version       int64      `json:"version"`
	Vin           []VinInfo  `json:"vin"`
	Vout          []VoutInfo `json:"vout"`
}

// AddressInfo is the address info for a returned address request
type AddressInfo struct {
	Address      string `json:"address"`
	IsMine       bool   `json:"ismine"`
	IsScript     bool   `json:"isscript"`
	IsValid      bool   `json:"isvalid"`
	IsWatchOnly  bool   `json:"iswatchonly"`
	ScriptPubKey string `json:"scriptPubKey"`
}

// AddressBalance is the address balance (unconfirmed and confirmed)
type AddressBalance struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
}

// AddressHistory is the history of transactions for an address
type AddressHistory []*HistoryRecord

// HistoryRecord is an internal record of AddressHistory
type HistoryRecord struct {
	Height int64  `json:"height"`
	TxHash string `json:"tx_hash"`
	TxPos  int64  `json:"tx_pos"`
	Value  int64  `json:"value"`
}

// Page is used as a sub-type for BlockInfo
type Page struct {
	Size int64    `json:"size"`
	URI  []string `json:"uri"`
}

// BlockPagesInfo is the response from the page request
type BlockPagesInfo []string

// CoinbaseTxInfo is the coinbase tx info inside the BlockInfo
type CoinbaseTxInfo struct {
	BlockHash     string     `json:"blockhash"`
	BlockTime     int64      `json:"blocktime"`
	Confirmations int64      `json:"confirmations"`
	Hash          string     `json:"hash"`
	Hex           string     `json:"hex"`
	LockTime      int64      `json:"locktime"`
	Size          int64      `json:"size"`
	Time          int64      `json:"time"`
	TxID          string     `json:"txid"`
	Version       int64      `json:"version"`
	Vin           []VinInfo  `json:"vin"`
	Vout          []VoutInfo `json:"vout"`
}

// VinInfo is the vin info inside the CoinbaseTxInfo
type VinInfo struct {
	Coinbase  string        `json:"coinbase"`
	ScriptSig ScriptSigInfo `json:"scriptSig"`
	Sequence  int64         `json:"sequence"`
	TxID      string        `json:"txid"`
	Vout      int64         `json:"vout"`
}

// ScriptSigInfo is the scriptSig info inside the VinInfo
type ScriptSigInfo struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// VinInfo is the vin info inside the CoinbaseTxInfo
type VoutInfo struct {
	N            int64            `json:"n"`
	ScriptPubKey ScriptPubKeyInfo `json:"scriptPubKey"`
	Value        float64          `json:"value"`
}

// ScriptPubKeyInfo is the scriptPubKey info inside the VoutInfo
type ScriptPubKeyInfo struct {
	Addresses []string `json:"addresses"`
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	OpReturn  string   `json:"opReturn"`
	ReqSigs   int64    `json:"reqSigs"`
	Type      string   `json:"type"`
}
