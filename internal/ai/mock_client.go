package ai

import "context"

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GenerateSummary(_ context.Context, prompt string) (string, error) {
	return "This wallet shows moderate transactional activity, repeated protocol interactions, and a medium risk profile based on observable on-chain behavior.", nil
}
