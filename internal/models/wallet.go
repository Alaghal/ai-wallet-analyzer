package models

type AnalyzeWalletRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
}

type AnalyzeWalletResponse struct {
	Address       string   `json:"address"`
	Chain         string   `json:"chain"`
	RiskScore     int      `json:"riskScore"`
	ActivityLevel string   `json:"activityLevel"`
	Tokens        []string `json:"tokens"`
	Summary       string   `json:"summary"`
}
