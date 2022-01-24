package whatsonchain

import "context"

// AddressService is the WhatsOnChain address related requests
type AddressService interface {
	AddressBalance(ctx context.Context, address string) (balance *AddressBalance, err error)
	AddressHistory(ctx context.Context, address string) (history AddressHistory, err error)
	AddressInfo(ctx context.Context, address string) (addressInfo *AddressInfo, err error)
	AddressUnspentTransactionDetails(ctx context.Context, address string, maxTransactions int) (history AddressHistory, err error)
	AddressUnspentTransactions(ctx context.Context, address string) (history AddressHistory, err error)
	BulkBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error)
}

// BlockService is the WhatsOnChain block related requests
type BlockService interface {
	GetBlockByHash(ctx context.Context, hash string) (blockInfo *BlockInfo, err error)
	GetBlockByHeight(ctx context.Context, height int64) (blockInfo *BlockInfo, err error)
	GetBlockPages(ctx context.Context, hash string, page int) (txList BlockPagesInfo, err error)
	GetHeaderByHash(ctx context.Context, hash string) (headerInfo *BlockInfo, err error)
	GetHeaders(ctx context.Context) (blockHeaders []*BlockInfo, err error)
}

// ChainService is the WhatsOnChain chain info requests
type ChainService interface {
	GetChainInfo(ctx context.Context) (chainInfo *ChainInfo, err error)
	GetCirculatingSupply(ctx context.Context) (supply float64, err error)
	GetExchangeRate(ctx context.Context) (rate *ExchangeRate, err error)
}

// DownloadService is the WhatsOnChain receipt and download related requests
type DownloadService interface {
	DownloadReceipt(ctx context.Context, hash string) (string, error)
	DownloadStatement(ctx context.Context, address string) (string, error)
}

// GeneralService is the WhatsOnChain general service requests
type GeneralService interface {
	GetExplorerLinks(ctx context.Context, query string) (results SearchResults, err error)
	GetHealth(ctx context.Context) (string, error)
}

// MempoolService is the WhatsOnChain mempool requests
type MempoolService interface {
	GetMempoolInfo(ctx context.Context) (info *MempoolInfo, err error)
	GetMempoolTransactions(ctx context.Context) (transactions []string, err error)
}

// ScriptService is the WhatsOnChain script requests
type ScriptService interface {
	BulkScriptUnspentTransactions(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error)
	GetScriptHistory(ctx context.Context, scriptHash string) (history ScriptList, err error)
	GetScriptUnspentTransactions(ctx context.Context, scriptHash string) (scriptList ScriptList, err error)
}

// TransactionService is the WhatsOnChain transaction related requests
type TransactionService interface {
	BroadcastTx(ctx context.Context, txHex string) (txID string, err error)
	BulkBroadcastTx(ctx context.Context, rawTxs []string, feedback bool) (response *BulkBroadcastResponse, err error)
	BulkTransactionDetails(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkTransactionDetailsProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkUnspentTransactions(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error)
	DecodeTransaction(ctx context.Context, txHex string) (txInfo *TxInfo, err error)
	GetMerkleProof(ctx context.Context, hash string) (merkleResults MerkleResults, err error)
	GetMerkleProofTSC(ctx context.Context, hash string) (merkleResults MerkleTSCResults, err error)
	GetRawTransactionData(ctx context.Context, hash string) (string, error)
	GetRawTransactionOutputData(ctx context.Context, hash string, vOutIndex int) (string, error)
	GetTxByHash(ctx context.Context, hash string) (txInfo *TxInfo, err error)
}

// ClientInterface is the WhatsOnChain client interface
type ClientInterface interface {
	AddressService
	BlockService
	ChainService
	DownloadService
	GeneralService
	MempoolService
	ScriptService
	TransactionService
	HTTPClient() HTTPInterface
	LastRequest() *LastRequest
	Network() NetworkType
	UserAgent() string
}
