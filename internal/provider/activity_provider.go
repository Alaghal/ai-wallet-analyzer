package provider

import (
	"context"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
)

type WalletActivityProvider interface {
	GetWalletActivity(ctx context.Context, address, chain string) (models.WalletActivity, error)
}
