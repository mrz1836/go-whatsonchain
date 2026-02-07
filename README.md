<div align="center">

# 🔗&nbsp;&nbsp;go-whatsonchain

**The unofficial Go SDK for the [whatsonchain.com API](https://docs.whatsonchain.com/) supporting both **[BSV](https://bsvblockchain.org/)** and **[BTC](https://thatsbtcnotbitcoin.com/)** blockchains.**

<br/>

<a href="https://github.com/mrz1836/go-whatsonchain/releases"><img src="https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/mrz1836/go-whatsonchain?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/mrz1836/go-whatsonchain/blob/master/LICENSE"><img src="https://img.shields.io/github/license/mrz1836/go-whatsonchain?style=flat-square&color=blue" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-whatsonchain/actions"><img src="https://img.shields.io/github/actions/workflow/status/mrz1836/go-whatsonchain/fortress.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/mrz1836/go-whatsonchain/actions"><img src="https://img.shields.io/github/last-commit/mrz1836/go-whatsonchain?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain"><img src="https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat-square" alt="Go Report"></a>
       <a href="https://codecov.io/gh/mrz1836/go-whatsonchain"><img src="https://codecov.io/gh/mrz1836/go-whatsonchain/branch/master/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/mrz1836/go-whatsonchain"><img src="https://api.scorecard.dev/projects/github.com/mrz1836/go-whatsonchain/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/mrz1836/go-whatsonchain/graphs/contributors"><img src="https://img.shields.io/github/contributors/mrz1836/go-whatsonchain?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://mrz1818.com/"><img src="https://img.shields.io/badge/donate-bitcoin-ff9900?style=flat-square&logo=bitcoin" alt="Bitcoin"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="50%">
       🚀&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="50%">
       💡&nbsp;<a href="#-usage"><code>Usage</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
    <td align="center">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
      ⚡&nbsp;<a href="#-benchmarks"><code>Benchmarks</code></a>
    </td>
    <td align="center">
      🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
      🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Usage</code></a>
    </td>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
    <td align="center">
       ⚖️&nbsp;<a href="#-license"><code>License</code></a>
    </td>
  </tr>
</table>

<br/>

## 📦 Installation

**go-whatsonchain** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get github.com/mrz1836/go-whatsonchain
```

<br/>

## 💡 Usage

### Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client with default options (BSV mainnet)
	client, err := whatsonchain.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded", client.UserAgent())
	log.Println("Chain:", client.Chain(), "Network:", client.Network())
}
```

### Configuration Options

The library uses functional options for clean and flexible configuration:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client with custom options
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBSV),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
		whatsonchain.WithAPIKey("your-secret-key"),
		whatsonchain.WithUserAgent("my-app/1.0"),
		whatsonchain.WithRateLimit(10),
		whatsonchain.WithRequestTimeout(60*time.Second),
		whatsonchain.WithRequestRetryCount(3),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded with custom options")
}
```

### Available Options

- `WithChain(chain)` - Set blockchain (ChainBSV or ChainBTC)
- `WithNetwork(network)` - Set network (NetworkMain, NetworkTest, NetworkStn)
- `WithAPIKey(key)` - Set API key for authenticated requests
- `WithUserAgent(agent)` - Set custom user agent
- `WithRateLimit(limit)` - Set rate limit per second
- `WithHTTPClient(client)` - Use custom HTTP client
- `WithRequestTimeout(timeout)` - Set request timeout
- `WithRequestRetryCount(count)` - Set retry count for failed requests
- `WithBackoff(initial, max, factor, jitter)` - Configure exponential backoff
- `WithDialer(keepAlive, timeout)` - Configure dialer settings
- `WithTransport(idle, tls, expect, maxIdle)` - Configure transport settings

### Multi-Chain Support

#### BSV Client

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create BSV client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBSV),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use BSV-specific methods
	opReturnData, err := client.GetOpReturnData(context.Background(), "your-tx-hash")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OP_RETURN data:", opReturnData)

	// Use shared methods (work for both BSV and BTC)
	chainInfo, err := client.GetChainInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BSV Chain Info: %+v", chainInfo)
}
```

#### BTC Client

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create BTC client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBTC),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use BTC-specific methods
	blockStats, err := client.GetBlockStats(context.Background(), 700000)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block Stats: %+v", blockStats)

	// Get miner statistics
	minerStats, err := client.GetMinerBlocksStats(context.Background(), 7) // last 7 days
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Miner Stats: %+v", minerStats)

	// Use shared methods (work for both BSV and BTC)
	chainInfo, err := client.GetChainInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BTC Chain Info: %+v", chainInfo)
}
```

<br/>

## 📚 Documentation
View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain?tab=doc)

<br/>

### Features
- **Multi-blockchain support** - Seamless switching between [BSV](https://bsvblockchain.org/) and [BTC](https://thatsbtcnotbitcoin.com/) blockchains with a single client
- **Production-ready HTTP client** - Built-in exponential backoff with configurable retry logic and crypto-secure jitter to handle transient failures gracefully
- **Intelligent rate limiting** - Per-second request throttling with automatic sleep intervals to stay within API quotas
- **Zero external dependencies** - Pure Go implementation with no production dependencies (testify only for testing)
- **Comprehensive API coverage** - 135+ endpoints (71 BSV, 64 BTC) fully implemented and tested
- **Flexible configuration** - Functional options pattern for clean, type-safe client initialization
- **Enterprise-grade transport** - Fine-grained control over timeouts, keep-alives, connection pooling, and TLS handshake settings
- **Network flexibility** - Switch between mainnet, testnet, and STN per client or per request

<br/>

> **Heads up!** `go-whatsonchain` is intentionally light on dependencies. The only
external package it uses is the excellent `testify` suite—and that's just for
our tests. You can drop this library into your projects without dragging along
extra baggage.

<br/>

<details>
<summary><strong><code>Supported API Coverage</code></strong></summary>
<br/>

**Coverage Summary:** 135 total endpoints (71 BSV + 64 BTC) from the [whatsonchain.com API](https://docs.whatsonchain.com/)

**Quick Navigation:** [BSV API](#bsv-api-71-endpoints) • [BTC API](#btc-api-64-endpoints) • [WebSockets](#websockets)

---

## BSV API (71 endpoints)

### ✅ Health (1 endpoint)
- [x] [Get Health Status](https://docs.whatsonchain.com/api/health) - `/woc`

### ✅ Chain Info (4 endpoints)
- [x] [Get Blockchain Info](https://docs.whatsonchain.com/api/chain-info#get-blockchain-info) - `/chain/info`
- [x] [Get Circulating Supply](https://docs.whatsonchain.com/api/chain-info#get-circulating-supply) - `/circulatingsupply`
- [x] [Get Chain Tips](https://docs.whatsonchain.com/api/chain-info#get-chain-tips) - `/chain/tips`
- [x] [Get Peer Info](https://docs.whatsonchain.com/api/chain-info#get-peer-info) - `/peer/info`

### ✅ Block (7 endpoints)
- [x] [Get Block by Hash](https://docs.whatsonchain.com/api/block#get-by-hash) - `/block/hash/{hash}`
- [x] [Get Block by Height](https://docs.whatsonchain.com/api/block#get-by-height) - `/block/height/{height}`
- [x] [Get Block Pages](https://docs.whatsonchain.com/api/block#get-block-pages) - `/block/hash/{hash}/page/{page}`
- [x] [Get Block Headers](https://docs.whatsonchain.com/api/block#get-headers) - `/block/headers`
- [x] [Get Header by Hash or Height](https://docs.whatsonchain.com/api/block#get-header-by-hash) - `/block/{hash}/header`
- [x] [Get Header Bytes File Links](https://docs.whatsonchain.com/api/block#get-header-bytes) - `/block/headers/resources`
- [x] [Get Latest Header Bytes](https://docs.whatsonchain.com/api/block#get-latest-headers) - `/block/headers/latest`

### ✅ Transaction (13 endpoints)
- [x] [Get Transaction by Hash](https://docs.whatsonchain.com/api/transaction#get-by-tx-hash) - `/tx/hash/{hash}`
- [x] [Get Transaction Propagation Status](https://docs.whatsonchain.com/api/transaction#get-tx-propagation) - `/tx/hash/{hash}/propagation`
- [x] [Broadcast Transaction](https://docs.whatsonchain.com/api/transaction#broadcast-transaction) - `/tx/raw` (POST)
- [x] [Bulk Transaction Details](https://docs.whatsonchain.com/api/transaction#bulk-transaction-details) - `/txs` (POST)
- [x] [Bulk Transaction Status](https://docs.whatsonchain.com/api/transaction#bulk-transaction-status) - `/txs/status` (POST)
- [x] [Decode Transaction](https://docs.whatsonchain.com/api/transaction#decode-transaction) - `/tx/decode` (POST)
- [x] [Download Receipt (PDF)](https://docs.whatsonchain.com/api/transaction#download-receipt) - Receipt download
- [x] [Get Transaction as Binary](https://docs.whatsonchain.com/api/transaction#get-tx-binary) - `/tx/{hash}/bin`
- [x] [Get Raw Transaction Data (Hex)](https://docs.whatsonchain.com/api/transaction#get-raw-tx-data) - `/tx/{hash}/hex`
- [x] [Bulk Raw Transaction Data](https://docs.whatsonchain.com/api/transaction#bulk-raw-tx-data) - `/txs/hex` (POST)
- [x] [Get Raw Transaction Output](https://docs.whatsonchain.com/api/transaction#get-raw-tx-output) - `/tx/{hash}/out/{index}/hex`
- [x] [Bulk Raw Transaction Output Data](https://docs.whatsonchain.com/api/transaction#bulk-raw-tx-output) - `/txs/vouts/hex` (POST)
- [x] [Get Merkle Proof (TSC)](https://docs.whatsonchain.com/api/transaction#get-merkle-proof) - `/tx/{hash}/proof/tsc`

### ✅ Mempool (2 endpoints)
- [x] [Get Mempool Info](https://docs.whatsonchain.com/api/mempool#get-mempool-info) - `/mempool/info`
- [x] [Get Mempool Transactions](https://docs.whatsonchain.com/api/mempool#get-mempool-transactions) - `/mempool/raw`

### ✅ (Un)Spent Transaction Outputs (14 endpoints)
- [x] [Get Unspent UTXOs by Address](https://docs.whatsonchain.com/api/address#get-unspent-transactions) - `/address/{address}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#get-unconfirmed-utxos) - `/address/{address}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-utxos) - `/addresses/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#get-confirmed-utxos) - `/address/{address}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#bulk-confirmed-utxos) - `/addresses/confirmed/unspent` (POST)
- [x] [Get Unspent UTXOs by Script](https://docs.whatsonchain.com/api/script#get-script-unspent-transactions) - `/script/{script}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#get-unconfirmed-script-utxos) - `/script/{script}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-utxos) - `/scripts/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#get-confirmed-script-utxos) - `/script/{script}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#bulk-confirmed-script-utxos) - `/scripts/confirmed/unspent` (POST)
- [x] [Get Unconfirmed Spent Output](https://docs.whatsonchain.com/api/utxo#get-unconfirmed-spent) - `/tx/{hash}/{index}/unconfirmed/spent`
- [x] [Get Confirmed Spent Output](https://docs.whatsonchain.com/api/utxo#get-confirmed-spent) - `/tx/{hash}/{index}/confirmed/spent`
- [x] [Get Spent Transaction Output](https://docs.whatsonchain.com/api/utxo#get-spent-output) - `/tx/{hash}/{index}/spent`
- [x] [Bulk Spent Transaction Outputs](https://docs.whatsonchain.com/api/utxo#bulk-spent-outputs) - `/utxos/spent` (POST)

### ✅ Address (13 endpoints)
- [x] [Get Address Info](https://docs.whatsonchain.com/api/address#get-address-info) - `/address/{address}/info`
- [x] [Get Address Usage Status](https://docs.whatsonchain.com/api/address#get-address-usage) - `/address/{address}/used`
- [x] [Get Associated Scripthashes](https://docs.whatsonchain.com/api/address#get-associated-scripthashes) - `/address/{address}/scripts`
- [x] [Download Statement (PDF)](https://docs.whatsonchain.com/api/address#download-statement) - Statement download
- [x] [Get Unconfirmed Balance](https://docs.whatsonchain.com/api/address#get-unconfirmed-balance) - `/address/{address}/unconfirmed/balance`
- [x] [Bulk Unconfirmed Balance](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-balance) - `/addresses/unconfirmed/balance` (POST)
- [x] [Get Confirmed Balance](https://docs.whatsonchain.com/api/address#get-confirmed-balance) - `/address/{address}/confirmed/balance`
- [x] [Bulk Confirmed Balance](https://docs.whatsonchain.com/api/address#bulk-confirmed-balance) - `/addresses/confirmed/balance` (POST)
- [x] [Get Unconfirmed History](https://docs.whatsonchain.com/api/address#get-unconfirmed-history) - `/address/{address}/unconfirmed/history`
- [x] [Bulk Unconfirmed History](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-history) - `/addresses/unconfirmed/history` (POST)
- [x] [Get Confirmed History](https://docs.whatsonchain.com/api/address#get-confirmed-history) - `/address/{address}/confirmed/history`
- [x] [Bulk Confirmed History](https://docs.whatsonchain.com/api/address#bulk-confirmed-history) - `/addresses/confirmed/history` (POST)
- [x] [Bulk History (All)](https://docs.whatsonchain.com/api/address#bulk-history) - `/addresses/history/all` (POST)

### ✅ Script (6 endpoints)
- [x] [Get Script Usage Status](https://docs.whatsonchain.com/api/script#get-script-usage) - `/script/{script}/used`
- [x] [Get Unconfirmed Script History](https://docs.whatsonchain.com/api/script#get-unconfirmed-script-history) - `/script/{script}/unconfirmed/history`
- [x] [Bulk Unconfirmed Script History](https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-history) - `/scripts/unconfirmed/history` (POST)
- [x] [Get Confirmed Script History](https://docs.whatsonchain.com/api/script#get-confirmed-script-history) - `/script/{script}/confirmed/history`
- [x] [Bulk Confirmed Script History](https://docs.whatsonchain.com/api/script#bulk-confirmed-script-history) - `/scripts/confirmed/history` (POST)

### ✅ Exchange Rate (2 endpoints)
- [x] [Get Current Exchange Rate](https://docs.whatsonchain.com/api/exchange-rate#get-exchange-rate) - `/exchangerate`
- [x] [Get Historical Exchange Rate](https://docs.whatsonchain.com/api/exchange-rate#get-historical-exchange-rate) - `/exchangerate/historical`

### ✅ Search (1 endpoint)
- [x] [Get Explorer Links](https://docs.whatsonchain.com/api/search#get-explorer-links) - `/search/links` (POST)

### ✅ On-Chain Data (1 endpoint)
- [x] [Get OP_RETURN Data](https://docs.whatsonchain.com/api/onchain-data#get-opreturn-data) - `/tx/{hash}/opreturn`

### ✅ Stats (6 endpoints)
- [x] [Get Block Stats by Height](https://docs.whatsonchain.com/api/stats#get-block-stats-by-height) - `/block/height/{height}/stats`
- [x] [Get Block Stats by Hash](https://docs.whatsonchain.com/api/stats#get-block-stats-by-hash) - `/block/hash/{hash}/stats`
- [x] [Get Miner Block Stats](https://docs.whatsonchain.com/api/stats#get-miner-block-stats) - `/miner/blocks/stats`
- [x] [Get Miner Minimum Fee Rate Stats](https://docs.whatsonchain.com/api/stats#get-miner-fee-stats) - `/miner/fees`
- [x] [Get Miner Summary Stats](https://docs.whatsonchain.com/api/stats#get-miner-summary-stats) - `/miner/summary/stats`
- [x] [Get Tag Count by Height](https://docs.whatsonchain.com/api/stats#get-tag-count-by-height) - `/block/tagcount/height/{height}/stats`

### ✅ Tokens (13 endpoints)

#### 1Sat Ordinals (7 endpoints)
- [x] [Get Token by Origin](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-origin) - `/token/1satordinals/{origin}/origin`
- [x] [Get Token by Outpoint](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-outpoint) - `/token/1satordinals/{outpoint}`
- [x] [Get Token Content](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-content) - `/token/1satordinals/{outpoint}/content`
- [x] [Get Token Latest Transfer](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-latest-transfer) - `/token/1satordinals/{outpoint}/latest`
- [x] [Get Token Transfer History](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-transfers-history) - `/token/1satordinals/{outpoint}/history`
- [x] [Get Tokens by TxID](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-tokens-by-txid) - `/token/1satordinals/tx/{txid}`
- [x] [Get 1Sat Ordinals Stats](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-stats) - `/tokens/1satordinals`

#### STAS v0 (6 endpoints)
- [x] [Get All STAS Tokens](https://docs.whatsonchain.com/api/tokens/stas#get-all-tokens) - `/tokens`
- [x] [Get STAS Token by ID](https://docs.whatsonchain.com/api/tokens/stas#get-token-by-id) - `/token/{contractId}/{symbol}`
- [x] [Get Token UTXOs for Address](https://docs.whatsonchain.com/api/tokens/stas#get-token-utxos-for-address) - `/address/{address}/tokens/unspent`
- [x] [Get Address Token Balance](https://docs.whatsonchain.com/api/tokens/stas#get-address-token-balance) - `/address/{address}/tokens`
- [x] [Get Token Transactions](https://docs.whatsonchain.com/api/tokens/stas#get-token-transactions) - `/token/{contractId}/{symbol}/tx`
- [x] [Get STAS Stats](https://docs.whatsonchain.com/api/tokens/stas#get-stats) - `/tokens/stas`

---

## BTC API (64 endpoints)

### ✅ Health (1 endpoint)
- [x] [Get Health Status](https://docs.whatsonchain.com/api/btc/health) - `/woc`

### ✅ Chain Info (4 endpoints)
- [x] [Get Blockchain Info](https://docs.whatsonchain.com/api/btc/chain-info#get-blockchain-info) - `/chain/info`
- [x] [Get Circulating Supply](https://docs.whatsonchain.com/api/btc/chain-info#get-circulating-supply) - `/circulatingsupply`
- [x] [Get Chain Tips](https://docs.whatsonchain.com/api/btc/chain-info#get-chain-tips) - `/chain/tips`
- [x] [Get Peer Info](https://docs.whatsonchain.com/api/btc/chain-info#get-peer-info) - `/peer/info`

### ✅ Block (7 endpoints)
- [x] [Get Block by Hash](https://docs.whatsonchain.com/api/btc/block#get-by-hash) - `/block/hash/{hash}`
- [x] [Get Block by Height](https://docs.whatsonchain.com/api/btc/block#get-by-height) - `/block/height/{height}`
- [x] [Get Block Pages](https://docs.whatsonchain.com/api/btc/block#get-block-pages) - `/block/hash/{hash}/page/{page}`
- [x] [Get Block Headers](https://docs.whatsonchain.com/api/btc/block#get-headers) - `/block/headers`
- [x] [Get Header by Hash or Height](https://docs.whatsonchain.com/api/btc/block#get-header-by-hash) - `/block/{hash}/header`
- [x] [Get Header Bytes File Links](https://docs.whatsonchain.com/api/btc/block#get-header-bytes) - `/block/headers/resources`
- [x] [Get Latest Header Bytes](https://docs.whatsonchain.com/api/btc/block#get-latest-headers) - `/block/headers/latest`

### ✅ Transaction (9 endpoints)
- [x] [Get Transaction by Hash](https://docs.whatsonchain.com/api/btc/transaction#get-by-tx-hash) - `/tx/hash/{hash}`
- [x] [Bulk Transaction Details](https://docs.whatsonchain.com/api/btc/transaction#bulk-transaction-details) - `/txs` (POST)
- [x] [Bulk Transaction Status](https://docs.whatsonchain.com/api/btc/transaction#bulk-transaction-status) - `/txs/status` (POST)
- [x] [Get Transaction as Binary](https://docs.whatsonchain.com/api/btc/transaction#get-tx-binary) - `/tx/{hash}/bin`
- [x] [Get Raw Transaction Data](https://docs.whatsonchain.com/api/btc/transaction#get-raw-tx-data) - `/tx/{hash}/hex`
- [x] [Bulk Raw Transaction Data](https://docs.whatsonchain.com/api/btc/transaction#bulk-raw-tx-data) - `/txs/hex` (POST)
- [x] [Get Raw Transaction Output](https://docs.whatsonchain.com/api/btc/transaction#get-raw-tx-output) - `/tx/{hash}/out/{index}/hex`
- [x] [Bulk Raw Transaction Output Data](https://docs.whatsonchain.com/api/btc/transaction#bulk-raw-tx-output) - `/txs/vouts/hex` (POST)
- [x] [Decode Transaction](https://docs.whatsonchain.com/api/btc/transaction#decode-transaction) - `/tx/decode` (POST)

### ✅ Mempool (2 endpoints)
- [x] [Get Mempool Info](https://docs.whatsonchain.com/api/btc/mempool#get-mempool-info) - `/mempool/info`
- [x] [Get Mempool Transactions](https://docs.whatsonchain.com/api/btc/mempool#get-mempool-transactions) - `/mempool/raw`

### ✅ (Un)Spent Transaction Outputs (14 endpoints)
- [x] [Get Unspent UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-unspent-transactions) - `/address/{address}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-utxos) - `/address/{address}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-utxos) - `/addresses/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-confirmed-utxos) - `/address/{address}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-utxos) - `/addresses/confirmed/unspent` (POST)
- [x] [Get Unspent UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-script-unspent-transactions) - `/script/{script}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-unconfirmed-script-utxos) - `/script/{script}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#bulk-unconfirmed-script-utxos) - `/scripts/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-confirmed-script-utxos) - `/script/{script}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#bulk-confirmed-script-utxos) - `/scripts/confirmed/unspent` (POST)
- [x] [Get Unconfirmed Spent Output](https://docs.whatsonchain.com/api/btc/utxo#get-unconfirmed-spent) - `/tx/{hash}/{index}/unconfirmed/spent`
- [x] [Get Confirmed Spent Output](https://docs.whatsonchain.com/api/btc/utxo#get-confirmed-spent) - `/tx/{hash}/{index}/confirmed/spent`
- [x] [Get Spent Transaction Output](https://docs.whatsonchain.com/api/btc/utxo#get-spent-output) - `/tx/{hash}/{index}/spent`
- [x] [Bulk Spent Transaction Outputs](https://docs.whatsonchain.com/api/btc/utxo#bulk-spent-outputs) - `/utxos/spent` (POST)

### ✅ Address (12 endpoints)
- [x] [Get Address Info](https://docs.whatsonchain.com/api/btc/address#get-address-info) - `/address/{address}/info`
- [x] [Get Address Usage Status](https://docs.whatsonchain.com/api/btc/address#get-address-usage) - `/address/{address}/used`
- [x] [Get Associated Scripthashes](https://docs.whatsonchain.com/api/btc/address#get-associated-scripthashes) - `/address/{address}/scripts`
- [x] [Get Unconfirmed Balance](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-balance) - `/address/{address}/unconfirmed/balance`
- [x] [Bulk Unconfirmed Balance](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-balance) - `/addresses/unconfirmed/balance` (POST)
- [x] [Get Confirmed Balance](https://docs.whatsonchain.com/api/btc/address#get-confirmed-balance) - `/address/{address}/confirmed/balance`
- [x] [Bulk Confirmed Balance](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-balance) - `/addresses/confirmed/balance` (POST)
- [x] [Get Unconfirmed History](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-history) - `/address/{address}/unconfirmed/history`
- [x] [Bulk Unconfirmed History](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-history) - `/addresses/unconfirmed/history` (POST)
- [x] [Get Confirmed History](https://docs.whatsonchain.com/api/btc/address#get-confirmed-history) - `/address/{address}/confirmed/history`
- [x] [Bulk Confirmed History](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-history) - `/addresses/confirmed/history` (POST)
- [x] [Bulk History (All)](https://docs.whatsonchain.com/api/btc/address#bulk-history) - `/addresses/history/all` (POST)

### ✅ Script (6 endpoints)
- [x] [Get Script Usage Status](https://docs.whatsonchain.com/api/btc/script#get-script-usage) - `/script/{script}/used`
- [x] [Get Unconfirmed Script History](https://docs.whatsonchain.com/api/btc/script#get-unconfirmed-script-history) - `/script/{script}/unconfirmed/history`
- [x] [Bulk Unconfirmed Script History](https://docs.whatsonchain.com/api/btc/script#bulk-unconfirmed-script-history) - `/scripts/unconfirmed/history` (POST)
- [x] [Get Confirmed Script History](https://docs.whatsonchain.com/api/btc/script#get-confirmed-script-history) - `/script/{script}/confirmed/history`
- [x] [Bulk Confirmed Script History](https://docs.whatsonchain.com/api/btc/script#bulk-confirmed-script-history) - `/scripts/confirmed/history` (POST)

### ✅ Exchange Rate (2 endpoints)
- [x] [Get Current Exchange Rate](https://docs.whatsonchain.com/api/btc/exchange-rate#get-exchange-rate) - `/exchangerate`
- [x] [Get Historical Exchange Rate](https://docs.whatsonchain.com/api/btc/exchange-rate#get-historical-exchange-rate) - `/exchangerate/historical`

### ✅ Search (1 endpoint)
- [x] [Get Explorer Links](https://docs.whatsonchain.com/api/btc/search#get-explorer-links) - `/search/links` (POST)

### ✅ Stats (6 endpoints)
- [x] [Get Block Stats by Height](https://docs.whatsonchain.com/api/btc/stats#get-block-stats-by-height) - `/block/height/{height}/stats`
- [x] [Get Block Stats by Hash](https://docs.whatsonchain.com/api/btc/stats#get-block-stats-by-hash) - `/block/hash/{hash}/stats`
- [x] [Get Miner Block Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-block-stats) - `/miner/blocks/stats`
- [x] [Get Miner Minimum Fee Rate Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-fee-stats) - `/miner/fees`
- [x] [Get Miner Summary Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-summary-stats) - `/miner/summary/stats`
- [x] [Get Tag Count by Height](https://docs.whatsonchain.com/api/btc/stats#get-tag-count-by-height) - `/block/tagcount/height/{height}/stats`

---

## WebSockets (ComingSoon™)
- [ ] [New block header event](https://docs.whatsonchain.com/api/websockets#new-block-header)
- [ ] [Block headers history](https://docs.whatsonchain.com/api/websockets#block-headers-history)
- [ ] [Block transactions](https://docs.whatsonchain.com/api/websockets#block-transactions)
- [ ] [Mempool transactions](https://docs.whatsonchain.com/api/websockets#mempool-transactions)
- [ ] [Confirmed transactions](https://docs.whatsonchain.com/api/websockets#confirmed-transactions)
- [ ] [Chain Stats](https://docs.whatsonchain.com/api/websockets#chain-stats)
- [ ] [Customized events](https://docs.whatsonchain.com/api/websockets#customized-events)
</details>

<details>
<summary><strong><code>Development Setup (Getting Started)</code></strong></summary>
<br/>

Install [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:

```bash
# Install MAGE-X for development and building
go install github.com/mrz1836/mage-x/cmd/magex@latest
magex update:install
```
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.

Then create and push a new Git tag using:

```bash
magex version:bump push=true bump=patch branch=master
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Build Commands</code></strong></summary>
<br/>

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong>GitHub Workflows</strong></summary>
<br/>

All workflows are driven by modular configuration in [`.github/env/`](.github/env/README.md) — no YAML editing required.

**[View all workflows and the control center →](.github/docs/workflows.md)**

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any managed tools. It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## 🧪 Examples & Tests

All unit tests and fuzz tests run via [GitHub Actions](https://github.com/mrz1836/go-whatsonchain/actions) and use [Go version 1.24.x](https://go.dev/doc/go1.24). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks

Run the Go benchmarks:

```bash script
magex bench time=2s
```

### Performance Results

Benchmarks run on **Apple M1 Max** using Go's built-in benchmark tool with 2-second intervals.

| Operation                          | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|------------------------------------|-------------:|--------------:|------------------------:|
| **Client Operations**              |
| Client Creation (Minimal)          |          275 |         1,048 |                       9 |
| Client Creation (Fully Configured) |          287 |         1,048 |                       9 |
| Build URL (Simple)                 |          158 |           144 |                       4 |
| Build URL (With Args)              |          194 |           168 |                       5 |
| Get Chain Config                   |          2.1 |             0 |                       0 |
| Set Chain Config                   |          2.1 |             0 |                       0 |
| **Address Operations**             |
| Get Address Info                   |        2,748 |         2,697 |                      25 |
| Get Address Balance                |        2,209 |         2,393 |                      24 |
| Get Address History                |        3,605 |         3,097 |                      32 |
| Get Address UTXOs                  |        5,668 |         3,705 |                      36 |
| Get Confirmed Balance              |        2,391 |         2,337 |                      22 |
| Get Unconfirmed Balance            |        2,683 |         2,337 |                      22 |
| **Transaction Operations**         |
| Get Transaction by Hash            |        8,117 |         6,035 |                      43 |
| Bulk Transaction Details (1 tx)    |       15,536 |        14,135 |                      69 |
| Bulk Transaction Details (20 txs)  |       17,489 |        16,815 |                      69 |
| Broadcast Transaction              |        2,071 |         2,753 |                      26 |
| Decode Transaction                 |        9,294 |         6,835 |                      54 |
| Get Merkle Proof                   |        5,679 |         3,826 |                      40 |
| Get Spent Output                   |        4,431 |         3,233 |                      31 |
| **Block Operations**               |
| Get Block by Hash                  |       11,271 |         6,307 |                      48 |
| Get Block by Height                |        8,691 |         5,963 |                      47 |
| Get Block Pages                    |        3,104 |         2,833 |                      30 |
| Get Header by Hash                 |       10,296 |         6,435 |                      49 |
| Get Latest Header Bytes            |        2,109 |         3,905 |                      21 |
| **Chain Info Operations**          |
| Get Chain Info                     |        3,442 |         2,881 |                      24 |
| Get Chain Tips                     |        2,546 |         2,497 |                      27 |
| Get Circulating Supply             |        1,514 |         1,993 |                      17 |
| Get Exchange Rate                  |        2,569 |         2,433 |                      26 |
| Get Historical Exchange Rate       |        3,916 |         3,089 |                      36 |
| Get Peer Info                      |        6,467 |         3,641 |                      33 |
| Get Mempool Info                   |        3,068 |         2,633 |                      27 |
| Get Mempool Transactions           |        2,648 |         2,593 |                      32 |

**Notes:**
- All times are in nanoseconds per operation
- Memory is bytes allocated per operation
- Benchmarks use mock HTTP responses for consistent, reproducible results

To reproduce these benchmarks:
```bash
magex bench time=2s
```

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with the AI assistants.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-whatsonchain&utm_term=go-whatsonchain&utm_content=go-whatsonchain) to ensure this journey continues indefinitely! :rocket:


[![Stars](https://img.shields.io/github/stars/mrz1836/go-whatsonchain?label=Please%20like%20us&style=social)](https://github.com/mrz1836/go-whatsonchain/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)](LICENSE)

