package whatsonchain

import "errors"

// ErrAddressNotFound is when an address is not found
var ErrAddressNotFound = errors.New("address not found")

// ErrBlockNotFound is when a block is not found
var ErrBlockNotFound = errors.New("block not found")

// ErrChainInfoNotFound is when the chain info is not found
var ErrChainInfoNotFound = errors.New("chain info not found")

// ErrExchangeRateNotFound is when the exchange rate is not found
var ErrExchangeRateNotFound = errors.New("exchange rate not found")

// ErrMempoolInfoNotFound is when the mempool info is not found
var ErrMempoolInfoNotFound = errors.New("mempool info not found")

// ErrHeadersNotFound is when the headers are not found
var ErrHeadersNotFound = errors.New("headers not found")

// ErrScriptNotFound is when a script is not found
var ErrScriptNotFound = errors.New("script not found")

// ErrTransactionNotFound is when a transaction is not found
var ErrTransactionNotFound = errors.New("transaction not found")

// ErrMaxAddressesExceeded is when the max addresses limit is exceeded
var ErrMaxAddressesExceeded = errors.New("max limit of addresses exceeded")

// ErrMaxScriptsExceeded is when the max scripts limit is exceeded
var ErrMaxScriptsExceeded = errors.New("max limit of scripts exceeded")

// ErrBroadcastFailed is when transaction broadcasting fails
var ErrBroadcastFailed = errors.New("error broadcasting transaction")

// ErrMaxTransactionsExceeded is when the max transactions limit is exceeded
var ErrMaxTransactionsExceeded = errors.New("max transactions limit exceeded")

// ErrMaxPayloadSizeExceeded is when the max payload size is exceeded
var ErrMaxPayloadSizeExceeded = errors.New("max overall payload size exceeded")

// ErrMaxTransactionSizeExceeded is when the max single transaction size is exceeded
var ErrMaxTransactionSizeExceeded = errors.New("max transaction size exceeded")

// ErrMaxUTXOsExceeded is when the max UTXO limit is exceeded
var ErrMaxUTXOsExceeded = errors.New("max limit of UTXOs exceeded")

// ErrMaxRawTransactionsExceeded is when the max raw transactions limit is exceeded
var ErrMaxRawTransactionsExceeded = errors.New("max limit of raw transactions exceeded")

// ErrMissingRequest is when a request is missing
var ErrMissingRequest = errors.New("missing request")

// ErrBadRequest is when a request is invalid
var ErrBadRequest = errors.New("bad request")
