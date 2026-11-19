# CLAUDE.md

Quick reference for Claude Code working on **go-whatsonchain** -- the unofficial Go SDK for the [whatsonchain.com API](https://docs.whatsonchain.com/).

## File Map

| File | Role |
|---|---|
| `client.go` | `Client` struct, `ClientOption` funcs, `clientOptions`, defaults |
| `whatsonchain.go` | `NewClient` constructor, `request()`, getters/setters |
| `interface.go` | All service interfaces + `ClientInterface` |
| `definitions.go` | Response types (`TxInfo`, `BlockInfo`, `AddressBalance`, etc.) |
| `errors.go` | Sentinel errors (`ErrTransactionNotFound`, `ErrInvalidNetwork`, etc.) |
| `url_builder.go` | `buildURL(path string, args ...any)` -- network-aware URL construction (BSV-only) |
| `request_helpers.go` | `requestAndUnmarshal[T]`, `requestAndUnmarshalSlice[T]`, `requestString`, `checkStatusCode`, `errorBodySnippet` |
| `http_client.go` | `RetryableHTTPClient` (retries, backoff, `Retry-After`), `ExponentialBackoff`, `SimpleHTTPClient` |
| `version.go` | `defaultUserAgent()` -- module version read from build info (`moduleVersion`) |
| `chunk.go` | `chunkSlice[T]` -- generic slice chunker for bulk processors |
| `addresses.go` | `AddressService` impl + `DownloadStatement` |
| `blocks.go` | `BlockService` impl |
| `transactions.go` | `TransactionService` impl + bulk processors + `GetOpReturnData` |
| `chain_info.go` | `ChainService` impl |
| `exchange_rates.go` | Exchange rate endpoints (part of `ChainService`) |
| `mempool.go` | `MempoolService` impl |
| `scripts.go` | `ScriptService` impl |
| `stats.go` | `StatsService` impl |
| `tokens.go` | `TokenService` impl (1Sat Ordinals, STAS) |
| `search.go` | `GetExplorerLinks` (part of `GeneralService`) |
| `health.go` | `GetHealth` (part of `GeneralService`) |
| `examples/` | Usage examples |

## Core Patterns

### 1. Functional Options (`client.go`)

```go
client, err := whatsonchain.NewClient(
    context.Background(),
    whatsonchain.WithNetwork(whatsonchain.NetworkMain),
    whatsonchain.WithAPIKey("key"),
    whatsonchain.WithRateLimit(10),
)
```

`ClientOption` is `func(*clientOptions)`. All `With*` functions modify the internal `clientOptions` struct. Auto-loads `WHATS_ON_CHAIN_API_KEY` env var if no explicit key, unless `WithoutEnvAPIKey()` is passed.

`WithHTTPClient` replaces the whole HTTP stack, so retries, backoff and timeouts don't apply; `WithRoundTripper` swaps only the transport and keeps them. Default backoff constants (`defaultBackOff*`) live in `client.go`.

### 2. Generic Request Helpers (`request_helpers.go`)

```go
// Single object -- returns *T
requestAndUnmarshal[T any](ctx, c, url, method, payload, emptyErr) (*T, error)

// Slice -- returns []T
requestAndUnmarshalSlice[T any](ctx, c, url, method, payload, emptyErr) ([]T, error)

// Raw string (GET only)
requestString(ctx, c, url, notFoundErr) (string, error)
```

`emptyErr` / `notFoundErr` is the endpoint's not-found sentinel (e.g., `ErrTransactionNotFound`). `checkStatusCode` maps a 404 to it: the bare sentinel for an empty body (how the API reports an unknown resource), or the sentinel wrapped with `HTTP 404: <body>` otherwise (e.g. `404 page not found` for an unknown route). `emptyErr` is also returned for an empty 200 body. Pass `nil` where a 404 is never "not found" (bulk POST, broadcast, health): the 404 then becomes `ErrRequestFailed`. Every other non-200 status returns `ErrRequestFailed: HTTP <code>: <body>`, with the body cut to 512 bytes by `errorBodySnippet`.

### 3. URL Building (`url_builder.go`)

```go
func (c *Client) buildURL(path string, args ...any) string
```

Constructs `https://api.whatsonchain.com/v1/bsv/{network}{path}`. String args are automatically `url.PathEscape`d.

### 4. Bulk Processors with `chunkSlice` + `time.Ticker`

Processors handle batching and rate limiting automatically. Pattern from `transactions.go:42`:

```go
batches := chunkSlice(hashes.TxIDs, MaxTransactionsUTXO)

ticker := time.NewTicker(time.Second / time.Duration(c.RateLimit()))
defer ticker.Stop()

for _, batch := range batches {
    select {
    case <-ctx.Done():
        return txList, ctx.Err()
    default:
    }
    select {
    case <-ctx.Done():
        return txList, ctx.Err()
    case <-ticker.C:
    }
    // process batch...
}
```

### 5. Service Interfaces (`interface.go`)

```go
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
    // + Getters/Setters
}
```

All services are implemented on `*Client`.

## Adding an Endpoint

1. **Add to interface** in `interface.go` under the appropriate service
2. **Define response type** in `definitions.go`
3. **Add error sentinel** to `errors.go` if needed
4. **Implement** in the service file:
   ```go
   func (c *Client) GetFoo(ctx context.Context, id string) (*FooResponse, error) {
       url := c.buildURL("/foo/%s", id)
       return requestAndUnmarshal[FooResponse](ctx, c, url, http.MethodGet, nil, ErrFooNotFound)
   }
   ```
5. **Write tests** -- table-driven with mock HTTP client

## Constants & Types (`definitions.go`)

**Chains:** `ChainBSV` ("bsv"). BSV is the only supported chain; `ChainType`, `ChainBSV`, `WithChain`, `SetChain`, and `Chain()` are retained as deprecated no-op shims for backward compatibility.

**Networks:** `NetworkMain` ("main"), `NetworkTest` ("test"), `NetworkStn` ("stn")

**Limits:**
| Constant | Value |
|---|---|
| `MaxTransactionsUTXO` | 20 |
| `MaxTransactionsRaw` | 20 |
| `MaxAddressesForLookup` | 20 |
| `MaxScriptsForLookup` | 20 |
| `MaxBroadcastTransactions` | 100 |
| `MaxSingleTransactionSize` | 102400 (100 KB) |
| `MaxCombinedTransactionSize` | 10,000,000 (10 MB) |

## Testing

### magex Commands
```bash
magex test            # Fast tests
magex test:race       # With race detector
magex test:cover      # Coverage report
magex bench           # Benchmarks
magex lint            # Linting (golangci-lint)
magex format:fix      # Format code
```

### Test Pattern
Tests use a mock HTTP client. Table-driven with `testCase` structs. `interface_test.go` verifies `Client` implements all service interfaces.

## Conventions

- Go 1.24+ (see `go.mod`)
- `context.Context` as first param on all API methods
- Pointer receivers on `*Client`
- Godoc on all exports with API doc link
- Custom errors from `errors.go` (prefix `Err`)
- Zero external deps (only `testify` for tests)
- Detailed standards in [AGENTS.md](AGENTS.md)

## Gotchas

1. **BSV-only**: BSV is the only supported chain. The chain-selection API (`WithChain`, `SetChain`, `Chain()`) is deprecated no-op shims; do not add new chain guards
2. **Max limits**: Bulk endpoints enforce limits; return wrapped errors (e.g., `ErrMaxTransactionsExceeded`)
3. **`emptyErr` / `notFoundErr` param**: Pass the endpoint's not-found sentinel to request helpers for 404/empty handling; pass `nil` for bulk POST endpoints (the API answers unknown items there with a 200 and per-item errors, so a 404 is a real failure)
4. **Rate limiting**: Bulk processors use `time.Ticker` -- never manual counter logic
5. **URL building**: Always use `c.buildURL()` -- it handles the `bsv/{network}` prefix + path escaping
6. **Concurrency**: `options` and `lastRequest` are mutex-protected; getters/setters are goroutine-safe
