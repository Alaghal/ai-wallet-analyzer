package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
)

// MockProvider - простой мок для WalletActivityProvider
type MockProvider struct {
	activity models.WalletActivity
	err      error
}

func (m *MockProvider) GetWalletActivity(ctx context.Context, address, chain string) (models.WalletActivity, error) {
	return m.activity, m.err
}

// MockAIClient - простой мок для ai.Client
type MockAIClient struct {
	summary string
	err     error
}

func (m *MockAIClient) GenerateSummary(ctx context.Context, prompt string) (string, error) {
	return m.summary, m.err
}

func TestCalculateRiskScore(t *testing.T) {
	tests := []struct {
		name     string
		activity models.WalletActivity
		want     int
	}{
		{
			name:     "Base score",
			activity: models.WalletActivity{},
			want:     20,
		},
		{
			name: "High transaction count",
			activity: models.WalletActivity{
				TransactionCount: 101,
			},
			want: 40,
		},
		{
			name: "Many unique interactions",
			activity: models.WalletActivity{
				UniqueInteractions: 11,
			},
			want: 35,
		},
		{
			name: "Many tokens",
			activity: models.WalletActivity{
				Tokens: []string{"T1", "T2", "T3", "T4", "T5", "T6"},
			},
			want: 30,
		},
		{
			name: "Combined high activity",
			activity: models.WalletActivity{
				TransactionCount:   150,
				UniqueInteractions: 15,
				Tokens:             []string{"T1", "T2", "T3", "T4", "T5", "T6"},
			},
			want: 65,
		},
		{
			name: "Max score limit",
			activity: models.WalletActivity{
				TransactionCount:   1000,
				UniqueInteractions: 100,
				Tokens:             make([]string, 50),
			},
			want: 65, // 20 + 20 + 15 + 10 = 65. Wait, current logic in analyzer.go is 20 + 20 + 15 + 10 = 65.
			// Let me re-read the code.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRiskScore(tt.activity); got != tt.want {
				t.Errorf("calculateRiskScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateActivityLevel(t *testing.T) {
	tests := []struct {
		count int
		want  string
	}{
		{0, "low"},
		{19, "low"},
		{20, "medium"},
		{99, "medium"},
		{100, "high"},
		{500, "high"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := calculateActivityLevel(tt.count); got != tt.want {
				t.Errorf("calculateActivityLevel(%d) = %v, want %v", tt.count, got, tt.want)
			}
		})
	}
}

func TestAnalyzerService_Analyze(t *testing.T) {
	ctx := context.Background()

	t.Run("Success with AI summary", func(t *testing.T) {
		provider := &MockProvider{
			activity: models.WalletActivity{
				Address:          "0x123",
				Chain:            "eth",
				TransactionCount: 50,
			},
		}
		aiClient := &MockAIClient{
			summary: "AI generated summary",
		}
		svc := NewAnalyzerService(provider, aiClient)

		res, err := svc.Analyze(ctx, "0x123", "eth")
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if res.Summary != "AI generated summary" {
			t.Errorf("expected AI summary, got %q", res.Summary)
		}
		if res.ActivityLevel != "medium" {
			t.Errorf("expected medium activity level, got %q", res.ActivityLevel)
		}
	})

	t.Run("Success with default summary when AI fails", func(t *testing.T) {
		provider := &MockProvider{
			activity: models.WalletActivity{
				Address:          "0x123",
				Chain:            "eth",
				TransactionCount: 5,
			},
		}
		aiClient := &MockAIClient{
			err: errors.New("ai error"),
		}
		svc := NewAnalyzerService(provider, aiClient)

		res, err := svc.Analyze(ctx, "0x123", "eth")
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		if res.Summary == "" || res.Summary == "AI generated summary" {
			t.Errorf("expected default summary, got %q", res.Summary)
		}
	})

	t.Run("Provider error", func(t *testing.T) {
		provider := &MockProvider{
			err: errors.New("provider error"),
		}
		svc := NewAnalyzerService(provider, nil)

		_, err := svc.Analyze(ctx, "0x123", "eth")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "provider error" {
			t.Errorf("expected 'provider error', got %v", err)
		}
	})
}
