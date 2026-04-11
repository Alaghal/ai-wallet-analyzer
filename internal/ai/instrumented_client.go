package ai

import (
	"context"
	"time"

	appMetrics "github.com/Alaghal/ai-wallet-analyzer/internal/metrics"
)

type InstrumentedClient struct {
	name    string
	client  Client
	metrics *appMetrics.Metrics
}

func NewInstrumentedClient(
	name string,
	client Client,
	metrics *appMetrics.Metrics,
) *InstrumentedClient {
	return &InstrumentedClient{
		name:    name,
		client:  client,
		metrics: metrics,
	}
}

func (c *InstrumentedClient) GenerateSummary(ctx context.Context, prompt string) (string, error) {
	start := time.Now()

	result, err := c.client.GenerateSummary(ctx, prompt)

	status := "success"
	if err != nil {
		status = "error"
	}

	c.metrics.LLMRequestsTotal.WithLabelValues(c.name, status).Inc()
	c.metrics.LLMRequestDuration.WithLabelValues(c.name, status).Observe(time.Since(start).Seconds())

	return result, err
}
