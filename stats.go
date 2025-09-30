package whatsonchain

import (
	"context"
	"net/http"
)

// BlockStats represents block statistics
type BlockStats struct {
	Height         int64   `json:"height"`
	Hash           string  `json:"hash"`
	Version        int     `json:"version"`
	Size           int     `json:"size"`
	Weight         int     `json:"weight"`
	MerkleRoot     string  `json:"merkleroot"`
	Timestamp      int64   `json:"timestamp"`
	MedianTime     int64   `json:"mediantime"`
	Nonce          int64   `json:"nonce"`
	Bits           string  `json:"bits"`
	Difficulty     float64 `json:"difficulty"`
	ChainWork      string  `json:"chainwork"`
	TxCount        int     `json:"tx_count"`
	TotalSize      int     `json:"total_size"`
	TotalFees      int64   `json:"total_fees"`
	SubsidyTotal   int64   `json:"subsidy_total"`
	SubsidyAddress int64   `json:"subsidy_address"`
	SubsidyMiner   int64   `json:"subsidy_miner"`
	MinerName      string  `json:"miner_name"`
	MinerAddress   string  `json:"miner_address"`
	FeeRateAvg     float64 `json:"fee_rate_avg"`
	FeeRateMin     float64 `json:"fee_rate_min"`
	FeeRateMax     float64 `json:"fee_rate_max"`
	FeeRateMedian  float64 `json:"fee_rate_median"`
	FeeRateStdDev  float64 `json:"fee_rate_stddev"`
	InputCount     int     `json:"input_count"`
	OutputCount    int     `json:"output_count"`
	UTXOIncrease   int     `json:"utxo_increase"`
	UTXOSizeInc    int     `json:"utxo_size_inc"`
}

// MinerStats represents miner statistics
type MinerStats struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	BlockCount int     `json:"block_count"`
	Percentage float64 `json:"percentage"`
}

// MinerFeeStats represents miner fee statistics
type MinerFeeStats struct {
	Timestamp int64   `json:"timestamp"`
	Name      string  `json:"name"`
	FeeRate   float64 `json:"fee_rate"`
}

// MinerSummaryStats represents miner summary statistics
type MinerSummaryStats struct {
	Days   int           `json:"days"`
	Miners []*MinerStats `json:"miners"`
}

// TagCount represents tag count statistics by height
type TagCount struct {
	Height    int64          `json:"height"`
	Hash      string         `json:"hash"`
	TagCounts map[string]int `json:"tag_counts"`
}

// GetBlockStats gets block statistics by height
//
// For more information: https://developers.whatsonchain.com/#block-stats
func (c *Client) GetBlockStats(ctx context.Context, height int64) (*BlockStats, error) {
	url := c.buildURL("/block/height/%d/stats", height)
	return requestAndUnmarshal[BlockStats](ctx, c, url, http.MethodGet, nil, nil)
}

// GetBlockStatsByHash gets block statistics by hash
//
// For more information: https://developers.whatsonchain.com/#block-stats
func (c *Client) GetBlockStatsByHash(ctx context.Context, hash string) (*BlockStats, error) {
	url := c.buildURL("/block/hash/%s/stats", hash)
	return requestAndUnmarshal[BlockStats](ctx, c, url, http.MethodGet, nil, nil)
}

// GetMinerBlocksStats gets miner blocks statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerBlocksStats(ctx context.Context, days int) ([]*MinerStats, error) {
	url := c.buildURL("/miner/blocks/stats?days=%d", days)
	return requestAndUnmarshalSlice[*MinerStats](ctx, c, url, http.MethodGet, nil, nil)
}

// GetMinerFeesStats gets miner fees statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerFeesStats(ctx context.Context, from, to int64) ([]*MinerFeeStats, error) {
	url := c.buildURL("/miner/fees?from=%d&to=%d", from, to)
	return requestAndUnmarshalSlice[*MinerFeeStats](ctx, c, url, http.MethodGet, nil, nil)
}

// GetMinerSummaryStats gets miner summary statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerSummaryStats(ctx context.Context, days int) (*MinerSummaryStats, error) {
	url := c.buildURL("/miner/summary/stats?days=%d", days)
	return requestAndUnmarshal[MinerSummaryStats](ctx, c, url, http.MethodGet, nil, nil)
}

// GetTagCountByHeight gets tag count statistics by height
//
// For more information: https://developers.whatsonchain.com/#tag-count-stats
func (c *Client) GetTagCountByHeight(ctx context.Context, height int64) (*TagCount, error) {
	url := c.buildURL("/block/tagcount/height/%d/stats", height)
	return requestAndUnmarshal[TagCount](ctx, c, url, http.MethodGet, nil, nil)
}
