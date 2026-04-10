package models

type AnalyzeWalletRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

type AnalyzeWalletResponse struct {
	Address            string   `json:"address"`
	Chain              string   `json:"chain"`
	RiskScore          int      `json:"riskScore"`
	ActivityLevel      string   `json:"activityLevel"`
	TransactionCount   int      `json:"transactionCount"`
	UniqueInteractions int      `json:"uniqueInteractions"`
	Tokens             []string `json:"tokens"`
	Summary            string   `json:"summary"`
}

type WalletActivity struct {
	Address            string
	Chain              string
	TransactionCount   int
	UniqueInteractions int
	Tokens             []string
}
