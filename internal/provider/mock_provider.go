package provider

import (
	"context"
	"strings"

	"github.com/yourname/ai-wallet-analyzer/internal/models"
)

type MockWalletActivityProvider struct{}

func NewMockWalletActivityProvider() *MockWalletActivityProvider {
	return &MockWalletActivityProvider{}
}

func (p *MockWalletActivityProvider) GetWalletActivity(
	_ context.Context,
	address, chain string,
) (models.WalletActivity, error) {
	normalizedChain := strings.TrimSpace(strings.ToLower(chain))
	if normalizedChain == "" {
		normalizedChain = "ethereum"
	}

	return models.WalletActivity{
		Address:            address,
		Chain:              normalizedChain,
		TransactionCount:   124,
		UniqueInteractions: 9,
		Tokens:             []string{"ETH", "USDT", "ARB"},
	}, nil
}
