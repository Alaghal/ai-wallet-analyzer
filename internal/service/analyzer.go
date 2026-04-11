package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Alaghal/ai-wallet-analyzer/internal/ai"
	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
	"github.com/Alaghal/ai-wallet-analyzer/internal/provider"
)

type AnalyzerService struct {
	provider  provider.WalletActivityProvider
	llmClient ai.Client
}

func NewAnalyzerService(
	provider provider.WalletActivityProvider,
	llmClient ai.Client,
) *AnalyzerService {
	return &AnalyzerService{
		provider:  provider,
		llmClient: llmClient,
	}
}

func (s *AnalyzerService) Analyze(
	ctx context.Context,
	address, chain string,
) (models.AnalyzeWalletResponse, error) {
	activity, err := s.provider.GetWalletActivity(ctx, address, chain)
	if err != nil {
		return models.AnalyzeWalletResponse{}, err
	}

	riskScore := calculateRiskScore(activity)
	activityLevel := calculateActivityLevel(activity.TransactionCount)

	summary := fmt.Sprintf(
		"Wallet %s on %s has %d transactions, %d unique interactions, and activity level %s.",
		activity.Address,
		activity.Chain,
		activity.TransactionCount,
		activity.UniqueInteractions,
		activityLevel,
	)

	if s.llmClient != nil {
		prompt := ai.BuildWalletSummaryPrompt(activity, riskScore, activityLevel)

		aiSummary, err := s.llmClient.GenerateSummary(ctx, prompt)
		if err != nil {
			log.Printf("llm summary generation failed: %v", err)
		} else if aiSummary != "" {
			summary = aiSummary
		}
	}

	return models.AnalyzeWalletResponse{
		Address:            activity.Address,
		Chain:              activity.Chain,
		RiskScore:          riskScore,
		ActivityLevel:      activityLevel,
		TransactionCount:   activity.TransactionCount,
		UniqueInteractions: activity.UniqueInteractions,
		Tokens:             activity.Tokens,
		Summary:            summary,
	}, nil
}

func calculateRiskScore(activity models.WalletActivity) int {
	score := 20

	if activity.TransactionCount > 100 {
		score += 20
	}
	if activity.UniqueInteractions > 10 {
		score += 15
	}
	if len(activity.Tokens) > 5 {
		score += 10
	}

	if score > 100 {
		score = 100
	}

	return score
}

func calculateActivityLevel(transactionCount int) string {
	switch {
	case transactionCount >= 100:
		return "high"
	case transactionCount >= 20:
		return "medium"
	default:
		return "low"
	}
}
