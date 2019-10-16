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
