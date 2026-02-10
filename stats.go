package whatsonchain

import (
	"context"
	"net/http"
)

// GetBlockStats gets block statistics by height
//
// For more information: https://developers.whatsonchain.com/#block-stats
func (c *Client) GetBlockStats(ctx context.Context, height int64) (*BlockStats, error) {
	url := c.buildURL("/block/height/%d/stats", height)
	return requestAndUnmarshal[BlockStats](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}

// GetBlockStatsByHash gets block statistics by hash
//
// For more information: https://developers.whatsonchain.com/#block-stats
func (c *Client) GetBlockStatsByHash(ctx context.Context, hash string) (*BlockStats, error) {
	url := c.buildURL("/block/hash/%s/stats", hash)
	return requestAndUnmarshal[BlockStats](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}

// GetMinerBlocksStats gets miner blocks statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerBlocksStats(ctx context.Context, days int) ([]*MinerStats, error) {
	url := c.buildURL("/miner/blocks/stats?days=%d", days)
	return requestAndUnmarshalSlice[*MinerStats](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}

// GetMinerFeesStats gets miner fees statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerFeesStats(ctx context.Context, from, to int64) ([]*MinerFeeStats, error) {
	url := c.buildURL("/miner/fees?from=%d&to=%d", from, to)
	return requestAndUnmarshalSlice[*MinerFeeStats](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}

// GetMinerSummaryStats gets miner summary statistics
//
// For more information: https://developers.whatsonchain.com/#miner-stats
func (c *Client) GetMinerSummaryStats(ctx context.Context, days int) (*MinerSummaryStats, error) {
	url := c.buildURL("/miner/summary/stats?days=%d", days)
	return requestAndUnmarshal[MinerSummaryStats](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}

// GetTagCountByHeight gets tag count statistics by height
//
// For more information: https://developers.whatsonchain.com/#tag-count-stats
func (c *Client) GetTagCountByHeight(ctx context.Context, height int64) (*TagCount, error) {
	url := c.buildURL("/block/tagcount/height/%d/stats", height)
	return requestAndUnmarshal[TagCount](ctx, c, url, http.MethodGet, nil, ErrStatsNotFound)
}
