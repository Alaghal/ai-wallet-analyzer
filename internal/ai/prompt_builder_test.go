package ai

import (
	"strings"
	"testing"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
)

func TestBuildWalletSummaryPrompt(t *testing.T) {
	activity := models.WalletActivity{
		Address:            "0x123",
		Chain:              "eth",
		TransactionCount:   50,
		UniqueInteractions: 5,
		Tokens:             []string{"USDC", "WETH"},
	}
	riskScore := 40
	activityLevel := "medium"

	prompt := BuildWalletSummaryPrompt(activity, riskScore, activityLevel)

	expectedParts := []string{
		"Wallet address: 0x123",
		"Chain: eth",
		"Transaction count: 50",
		"Unique interactions: 5",
		"Detected tokens: USDC, WETH",
		"Calculated risk score: 40",
		"Calculated activity level: medium",
	}

	for _, part := range expectedParts {
		if !strings.Contains(prompt, part) {
			t.Errorf("prompt does not contain expected part: %q", part)
		}
	}
}
