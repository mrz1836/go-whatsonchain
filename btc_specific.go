package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// BTCService is the interface for BTC-specific endpoints
type BTCService interface {
	GetBlockStats(ctx context.Context, height int64) (*BlockStats, error)
	GetBlockStatsByHash(ctx context.Context, hash string) (*BlockStats, error)
	GetMinerBlocksStats(ctx context.Context, days int) ([]*MinerStats, error)
	GetMinerFeesStats(ctx context.Context, from, to int64) ([]*MinerFeeStats, error)
	GetMinerSummaryStats(ctx context.Context, days int) (*MinerSummaryStats, error)
}

// BlockStats represents block statistics (BTC-specific)
type BlockStats struct {
	AvgFee        int64   `json:"avgfee"`
	AvgFeeRate    float64 `json:"avgfeerate"`
	AvgTxSize     float64 `json:"avgtxsize"`
	BlockHash     string  `json:"blockhash"`
	FeeRatePerc   []int   `json:"feerate_percentiles"`
	Height        int64   `json:"height"`
	Ins           int64   `json:"ins"`
	MaxFee        int64   `json:"maxfee"`
	MaxFeeRate    float64 `json:"maxfeerate"`
	MaxTxSize     int64   `json:"maxtxsize"`
	MedianFee     int64   `json:"medianfee"`
	MedianTime    int64   `json:"mediantime"`
	MedianTxSize  int64   `json:"mediantxsize"`
	MinFee        int64   `json:"minfee"`
	MinFeeRate    float64 `json:"minfeerate"`
	MinTxSize     int64   `json:"mintxsize"`
	Outs          int64   `json:"outs"`
	Subsidy       int64   `json:"subsidy"`
	SWTotalSize   int64   `json:"swtotal_size"`
	SWTotalWeight int64   `json:"swtotal_weight"`
	SWTxs         int64   `json:"swtxs"`
	Time          int64   `json:"time"`
	TotalOut      int64   `json:"total_out"`
	TotalSize     int64   `json:"total_size"`
	TotalWeight   int64   `json:"total_weight"`
	TotalFee      int64   `json:"totalfee"`
	Txs           int64   `json:"txs"`
	UTXOIncrease  int64   `json:"utxo_increase"`
	UTXOSizeInc   int64   `json:"utxo_size_inc"`
}

// MinerStats represents miner statistics (BTC-specific)
type MinerStats struct {
	Miner      string  `json:"miner"`
	BlockCount int64   `json:"block_count"`
	Percentage float64 `json:"percentage"`
}

// MinerFeeStats represents miner fee statistics (BTC-specific)
type MinerFeeStats struct {
	Miner    string  `json:"miner"`
	TotalFee int64   `json:"total_fee"`
	AvgFee   float64 `json:"avg_fee"`
}

// MinerSummaryStats represents miner summary statistics (BTC-specific)
type MinerSummaryStats struct {
	TotalBlocks int64         `json:"total_blocks"`
	Miners      []*MinerStats `json:"miners"`
}

// GetBlockStats gets block statistics by height (BTC-only endpoint)
//
// For more information: https://developers.whatsonchain.com/#btc-block-stats
func (c *Client) GetBlockStats(ctx context.Context, height int64) (*BlockStats, error) {
	// Only available for BTC
	if c.Chain() != ChainBTC {
		return nil, ErrBTCChainRequired
	}

	var resp string
	var err error
	// https://api.whatsonchain.com/v1/btc/<network>/block/height/<height>/stats
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/height/%d/stats", apiEndpointBase, c.Chain(), c.Network(), height),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var stats *BlockStats
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &stats)
	}
	return stats, err
}

// GetBlockStatsByHash gets block statistics by hash (BTC-only endpoint)
//
// For more information: https://developers.whatsonchain.com/#btc-block-stats
func (c *Client) GetBlockStatsByHash(ctx context.Context, hash string) (*BlockStats, error) {
	// Only available for BTC
	if c.Chain() != ChainBTC {
		return nil, ErrBTCChainRequired
	}

	var resp string
	var err error
	// https://api.whatsonchain.com/v1/btc/<network>/block/hash/<hash>/stats
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/hash/%s/stats", apiEndpointBase, c.Chain(), c.Network(), hash),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var stats *BlockStats
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &stats)
	}
	return stats, err
}

// GetMinerBlocksStats gets miner blocks statistics (BTC-only endpoint)
//
// For more information: https://developers.whatsonchain.com/#btc-miner-stats
func (c *Client) GetMinerBlocksStats(ctx context.Context, days int) ([]*MinerStats, error) {
	// Only available for BTC
	if c.Chain() != ChainBTC {
		return nil, ErrBTCChainRequired
	}

	var resp string
	var err error
	// https://api.whatsonchain.com/v1/btc/<network>/miner/blocks/stats?days=<days>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/miner/blocks/stats?days=%d", apiEndpointBase, c.Chain(), c.Network(), days),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var stats []*MinerStats
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &stats)
	}
	return stats, err
}

// GetMinerFeesStats gets miner fees statistics (BTC-only endpoint)
//
// For more information: https://developers.whatsonchain.com/#btc-miner-stats
func (c *Client) GetMinerFeesStats(ctx context.Context, from, to int64) ([]*MinerFeeStats, error) {
	// Only available for BTC
	if c.Chain() != ChainBTC {
		return nil, ErrBTCChainRequired
	}

	var resp string
	var err error
	// https://api.whatsonchain.com/v1/btc/<network>/miner/fees?from=<from>&to=<to>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/miner/fees?from=%d&to=%d", apiEndpointBase, c.Chain(), c.Network(), from, to),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var stats []*MinerFeeStats
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &stats)
	}
	return stats, err
}

// GetMinerSummaryStats gets miner summary statistics (BTC-only endpoint)
//
// For more information: https://developers.whatsonchain.com/#btc-miner-stats
func (c *Client) GetMinerSummaryStats(ctx context.Context, days int) (*MinerSummaryStats, error) {
	// Only available for BTC
	if c.Chain() != ChainBTC {
		return nil, ErrBTCChainRequired
	}

	var resp string
	var err error
	// https://api.whatsonchain.com/v1/btc/<network>/miner/summary/stats?days=<days>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/miner/summary/stats?days=%d", apiEndpointBase, c.Chain(), c.Network(), days),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var stats *MinerSummaryStats
	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &stats)
	}
	return stats, err
}
