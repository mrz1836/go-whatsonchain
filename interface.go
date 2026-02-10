package whatsonchain

import (
	"context"
	"time"
)

// AddressService is the WhatsOnChain address related requests
type AddressService interface {
	// Deprecated: AddressBalance uses a combined endpoint no longer in the API. Use AddressConfirmedBalance and AddressUnconfirmedBalance.
	AddressBalance(ctx context.Context, address string) (balance *AddressBalance, err error)
	AddressConfirmedBalance(ctx context.Context, address string) (balance *AddressConfirmedBalance, err error)
	AddressConfirmedHistory(ctx context.Context, address string) (history AddressHistory, err error)
	AddressConfirmedUTXOs(ctx context.Context, address string) (history AddressHistory, err error)
	// Deprecated: AddressHistory uses a combined endpoint no longer in the API. Use AddressConfirmedHistory and AddressUnconfirmedHistory.
	AddressHistory(ctx context.Context, address string) (history AddressHistory, err error)
	AddressInfo(ctx context.Context, address string) (addressInfo *AddressInfo, err error)
	AddressScripts(ctx context.Context, address string) (scripts *AddressScripts, err error)
	AddressUnconfirmedBalance(ctx context.Context, address string) (balance *AddressUnconfirmedBalance, err error)
	AddressUnconfirmedHistory(ctx context.Context, address string) (history AddressHistory, err error)
	AddressUnconfirmedUTXOs(ctx context.Context, address string) (history AddressHistory, err error)
	AddressUnspentTransactionDetails(ctx context.Context, address string, maxTransactions int) (history AddressHistory, err error)
	AddressUnspentTransactions(ctx context.Context, address string) (history AddressHistory, err error)
	AddressUsed(ctx context.Context, address string) (used *AddressUsed, err error)
	BulkAddressConfirmedBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error)
	BulkAddressConfirmedHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error)
	BulkAddressConfirmedUTXOs(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error)
	BulkAddressHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error)
	BulkAddressUnconfirmedBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error)
	BulkAddressUnconfirmedHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error)
	BulkAddressUnconfirmedUTXOs(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error)
	// Deprecated: BulkBalance uses a combined endpoint no longer in the API. Use BulkAddressConfirmedBalance and BulkAddressUnconfirmedBalance.
	BulkBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error)
}

// BlockService is the WhatsOnChain block related requests
type BlockService interface {
	GetBlockByHash(ctx context.Context, hash string) (blockInfo *BlockInfo, err error)
	GetBlockByHeight(ctx context.Context, height int64) (blockInfo *BlockInfo, err error)
	GetBlockPages(ctx context.Context, hash string, page int) (txList BlockPagesInfo, err error)
	GetHeaderByHash(ctx context.Context, hash string) (headerInfo *BlockInfo, err error)
	GetHeaders(ctx context.Context) (blockHeaders []*BlockInfo, err error)
	GetHeaderBytesFileLinks(ctx context.Context) (resource *HeaderBytesResource, err error)
	GetLatestHeaderBytes(ctx context.Context, count int) (headerBytes string, err error)
}

// ChainService is the WhatsOnChain chain info requests
type ChainService interface {
	GetChainInfo(ctx context.Context) (chainInfo *ChainInfo, err error)
	GetChainTips(ctx context.Context) (chainTips []*ChainTip, err error)
	GetCirculatingSupply(ctx context.Context) (supply float64, err error)
	GetExchangeRate(ctx context.Context) (rate *ExchangeRate, err error)
	GetHistoricalExchangeRate(ctx context.Context, from, to int64) (rates []*HistoricalExchangeRate, err error)
	GetPeerInfo(ctx context.Context) (peerInfo []*PeerInfo, err error)
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
	BulkScriptConfirmedHistory(ctx context.Context, list *ScriptsList) (response BulkScriptHistoryResponse, err error)
	BulkScriptConfirmedUTXOs(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error)
	BulkScriptUnconfirmedHistory(ctx context.Context, list *ScriptsList) (response BulkScriptHistoryResponse, err error)
	BulkScriptUnconfirmedUTXOs(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error)
	// Deprecated: BulkScriptUnspentTransactions uses a combined endpoint no longer in the API. Use BulkScriptConfirmedUTXOs and BulkScriptUnconfirmedUTXOs.
	BulkScriptUnspentTransactions(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error)
	GetScriptConfirmedHistory(ctx context.Context, scriptHash string) (history ScriptList, err error)
	// Deprecated: GetScriptHistory uses a combined endpoint no longer in the API. Use GetScriptConfirmedHistory and GetScriptUnconfirmedHistory.
	GetScriptHistory(ctx context.Context, scriptHash string) (history ScriptList, err error)
	GetScriptUnconfirmedHistory(ctx context.Context, scriptHash string) (history ScriptList, err error)
	GetScriptUnspentTransactions(ctx context.Context, scriptHash string) (scriptList ScriptList, err error)
	GetScriptUsed(ctx context.Context, scriptHash string) (used bool, err error)
	ScriptConfirmedUTXOs(ctx context.Context, scriptHash string) (scriptList ScriptList, err error)
	ScriptUnconfirmedUTXOs(ctx context.Context, scriptHash string) (scriptList ScriptList, err error)
}

// StatsService is the WhatsOnChain stats requests
type StatsService interface {
	GetBlockStats(ctx context.Context, height int64) (*BlockStats, error)
	GetBlockStatsByHash(ctx context.Context, hash string) (*BlockStats, error)
	GetMinerBlocksStats(ctx context.Context, days int) ([]*MinerStats, error)
	GetMinerFeesStats(ctx context.Context, from, to int64) ([]*MinerFeeStats, error)
	GetMinerSummaryStats(ctx context.Context, days int) (*MinerSummaryStats, error)
	GetTagCountByHeight(ctx context.Context, height int64) (*TagCount, error)
}

// TokenService is the WhatsOnChain token related requests
type TokenService interface {
	// 1Sat Ordinals
	GetOneSatOrdinalByOrigin(ctx context.Context, origin string) (token *OneSatOrdinalToken, err error)
	GetOneSatOrdinalByOutpoint(ctx context.Context, outpoint string) (token *OneSatOrdinalToken, err error)
	GetOneSatOrdinalContent(ctx context.Context, outpoint string) (content *OneSatOrdinalContent, err error)
	GetOneSatOrdinalLatest(ctx context.Context, outpoint string) (latest *OneSatOrdinalLatest, err error)
	GetOneSatOrdinalHistory(ctx context.Context, outpoint string) (history []*OneSatOrdinalHistory, err error)
	GetOneSatOrdinalsByTxID(ctx context.Context, txid string) (tokens []*OneSatOrdinalToken, err error)
	GetOneSatOrdinalsStats(ctx context.Context) (stats *OneSatOrdinalStats, err error)

	// STAS v0
	GetAllSTASTokens(ctx context.Context) (tokens []*STASToken, err error)
	GetSTASTokenByID(ctx context.Context, contractID, symbol string) (token *STASToken, err error)
	GetTokenUTXOsForAddress(ctx context.Context, address string) (utxos []*STASTokenUTXO, err error)
	GetAddressTokenBalance(ctx context.Context, address string) (balance *STASTokenBalance, err error)
	GetTokenTransactions(ctx context.Context, contractID, symbol string) (transactions TxList, err error)
	GetSTASStats(ctx context.Context) (stats *STASStats, err error)
}

// TransactionService is the WhatsOnChain transaction related requests
type TransactionService interface {
	BroadcastTx(ctx context.Context, txHex string) (txID string, err error)
	// Deprecated: BulkBroadcastTx uses an endpoint no longer in the API. Use BroadcastTx instead.
	BulkBroadcastTx(ctx context.Context, rawTxs []string, feedback bool) (response *BulkBroadcastResponse, err error)
	BulkRawTransactionData(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkRawTransactionDataProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkRawTransactionOutputData(ctx context.Context, request *BulkRawOutputRequest) (responses []*BulkRawOutputResponse, err error)
	BulkSpentOutputs(ctx context.Context, request *BulkSpentOutputRequest) (response BulkSpentOutputResponse, err error)
	BulkTransactionDetails(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkTransactionDetailsProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error)
	BulkTransactionStatus(ctx context.Context, hashes *TxHashes) (txStatusList TxStatusList, err error)
	// Deprecated: BulkUnspentTransactions uses an endpoint no longer in the API. Use BulkAddressConfirmedUTXOs and BulkAddressUnconfirmedUTXOs.
	BulkUnspentTransactions(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error)
	// Deprecated: BulkUnspentTransactionsProcessor wraps BulkUnspentTransactions which is deprecated.
	BulkUnspentTransactionsProcessor(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error)
	DecodeTransaction(ctx context.Context, txHex string) (txInfo *TxInfo, err error)
	GetConfirmedSpentOutput(ctx context.Context, txHash string, index int) (spentOutput *SpentOutput, err error)
	// Deprecated: GetMerkleProof uses a non-TSC endpoint no longer in the API. Use GetMerkleProofTSC instead.
	GetMerkleProof(ctx context.Context, hash string) (merkleResults MerkleResults, err error)
	GetMerkleProofTSC(ctx context.Context, hash string) (merkleResults MerkleTSCResults, err error)
	GetRawTransactionData(ctx context.Context, hash string) (string, error)
	GetRawTransactionOutputData(ctx context.Context, hash string, vOutIndex int) (string, error)
	GetSpentOutput(ctx context.Context, txHash string, index int) (spentOutput *SpentOutput, err error)
	GetTransactionAsBinary(ctx context.Context, hash string) ([]byte, error)
	GetTransactionPropagationStatus(ctx context.Context, hash string) (propagationStatus *PropagationStatus, err error)
	GetTxByHash(ctx context.Context, hash string) (txInfo *TxInfo, err error)
	GetUnconfirmedSpentOutput(ctx context.Context, txHash string, index int) (spentOutput *SpentOutput, err error)
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
	StatsService
	TokenService
	TransactionService
	BSVService
	BTCService

	// Getters
	APIKey() string
	BackoffConfig() (initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration)
	Chain() ChainType
	DialerConfig() (keepAlive, timeout time.Duration)
	HTTPClient() HTTPInterface
	LastRequest() *LastRequest
	Network() NetworkType
	RateLimit() int
	RequestRetryCount() int
	RequestTimeout() time.Duration
	TransportConfig() (idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int)
	UserAgent() string

	// Setters
	SetAPIKey(apiKey string)
	SetBackoffConfig(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitter time.Duration)
	SetChain(chain ChainType)
	SetDialerConfig(keepAlive, timeout time.Duration)
	SetNetwork(network NetworkType)
	SetRateLimit(rateLimit int)
	SetRequestRetryCount(count int)
	SetRequestTimeout(timeout time.Duration)
	SetTransportConfig(idleTimeout, tlsTimeout, expectContinueTimeout time.Duration, maxIdleConnections int)
	SetUserAgent(userAgent string)
}
