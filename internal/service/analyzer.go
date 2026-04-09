package service

import (
	"fmt"
	"strings"

	"github.com/yourname/ai-wallet-analyzer/internal/models"
)

type AnalyzerService struct{}

func NewAnalyzerService() *AnalyzerService {
	return &AnalyzerService{}
}

func (s *AnalyzerService) Analyze(address, chain string) models.AnalyzeWalletResponse {
	normalizedChain := strings.TrimSpace(strings.ToLower(chain))
	if normalizedChain == "" {
		normalizedChain = "ethereum"
	}

	return models.AnalyzeWalletResponse{
		Address:       address,
		Chain:         normalizedChain,
		RiskScore:     42,
		ActivityLevel: "medium",
		Tokens:        []string{"ETH", "USDT", "ARB"},
		Summary: fmt.Sprintf(
			"Wallet %s on %s shows medium activity with interactions across common tokens and DeFi-related behavior.",
			address,
			normalizedChain,
		),
	}
}
