package ai

import "context"

type Client interface {
	GenerateSummary(ctx context.Context, prompt string) (string, error)
}
