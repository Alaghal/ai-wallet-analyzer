package ai

import (
	"fmt"
	"strings"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
)

func BuildWalletSummaryPrompt(activity models.WalletActivity, riskScore int, activityLevel string) string {
	tokens := "none"
	if len(activity.Tokens) > 0 {
		tokens = strings.Join(activity.Tokens, ", ")
	}

	return fmt.Sprintf(`
Analyze the following Web3 wallet activity and write a short professional summary.

Wallet address: %s
Chain: %s
Transaction count: %d
Unique interactions: %d
Detected tokens: %s
Calculated risk score: %d
Calculated activity level: %s

Return a concise summary in 2-3 sentences.
Focus on activity style, likely behavior, and risk signals.
Do not mention that this was generated from a prompt.
`, activity.Address, activity.Chain, activity.TransactionCount, activity.UniqueInteractions, tokens, riskScore, activityLevel)
}
