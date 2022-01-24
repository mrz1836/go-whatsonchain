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
