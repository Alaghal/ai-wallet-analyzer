package provider

import (
	"context"
	"time"

	appMetrics "github.com/Alaghal/ai-wallet-analyzer/internal/metrics"
	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
)

type InstrumentedProvider struct {
	name     string
	provider WalletActivityProvider
	metrics  *appMetrics.Metrics
}

func NewInstrumentedProvider(
	name string,
	provider WalletActivityProvider,
	metrics *appMetrics.Metrics,
) *InstrumentedProvider {
	return &InstrumentedProvider{
		name:     name,
		provider: provider,
		metrics:  metrics,
	}
}

func (p *InstrumentedProvider) GetWalletActivity(
	ctx context.Context,
	address, chain string,
) (models.WalletActivity, error) {
	start := time.Now()

	result, err := p.provider.GetWalletActivity(ctx, address, chain)

	status := "success"
	if err != nil {
		status = "error"
	}

	p.metrics.ProviderRequestsTotal.WithLabelValues(p.name, status).Inc()
	p.metrics.ProviderRequestDuration.WithLabelValues(p.name, status).Observe(time.Since(start).Seconds())

	return result, err
}
