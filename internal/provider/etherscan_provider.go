package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourname/ai-wallet-analyzer/internal/models"
)

type EtherscanWalletActivityProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewEtherscanWalletActivityProvider(
	baseURL string,
	apiKey string,
	timeout time.Duration,
) *EtherscanWalletActivityProvider {
	return &EtherscanWalletActivityProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type etherscanTxListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		From        string `json:"from"`
		To          string `json:"to"`
		Contract    string `json:"contractAddress"`
		TokenSymbol string `json:"tokenSymbol"`
	} `json:"result"`
}

func (p *EtherscanWalletActivityProvider) GetWalletActivity(
	ctx context.Context,
	address, chain string,
) (models.WalletActivity, error) {
	normalizedChain := strings.TrimSpace(strings.ToLower(chain))
	if normalizedChain == "" {
		normalizedChain = "ethereum"
	}

	if normalizedChain != "ethereum" {
		return models.WalletActivity{}, fmt.Errorf("unsupported chain: %s", normalizedChain)
	}

	u, err := url.Parse(p.baseURL)
	if err != nil {
		return models.WalletActivity{}, fmt.Errorf("parse provider url: %w", err)
	}

	query := u.Query()
	query.Set("module", "account")
	query.Set("action", "txlist")
	query.Set("address", address)
	query.Set("startblock", "0")
	query.Set("endblock", "99999999")
	query.Set("page", "1")
	query.Set("offset", "50")
	query.Set("sort", "desc")
	if p.apiKey != "" {
		query.Set("apikey", p.apiKey)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.WalletActivity{}, fmt.Errorf("create provider request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return models.WalletActivity{}, fmt.Errorf("provider request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.WalletActivity{}, fmt.Errorf("provider returned status %d", resp.StatusCode)
	}

	var payload etherscanTxListResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return models.WalletActivity{}, fmt.Errorf("decode provider response: %w", err)
	}

	interactionSet := make(map[string]struct{})
	tokenSet := make(map[string]struct{})

	for _, tx := range payload.Result {
		if tx.To != "" {
			interactionSet[strings.ToLower(tx.To)] = struct{}{}
		}
		if tx.TokenSymbol != "" {
			tokenSet[strings.ToUpper(tx.TokenSymbol)] = struct{}{}
		}
	}

	tokens := make([]string, 0, len(tokenSet))
	for token := range tokenSet {
		tokens = append(tokens, token)
	}

	return models.WalletActivity{
		Address:            address,
		Chain:              normalizedChain,
		TransactionCount:   len(payload.Result),
		UniqueInteractions: len(interactionSet),
		Tokens:             tokens,
	}, nil
}
