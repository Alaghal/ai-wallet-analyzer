package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
	"github.com/Alaghal/ai-wallet-analyzer/internal/service"
)

type MockProvider struct {
	activity models.WalletActivity
	err      error
}

func (m *MockProvider) GetWalletActivity(ctx context.Context, address, chain string) (models.WalletActivity, error) {
	return m.activity, m.err
}

func TestWalletHandler_AnalyzeWallet(t *testing.T) {
	provider := &MockProvider{}
	svc := service.NewAnalyzerService(provider, nil)
	handler := NewWalletHandler(svc)

	t.Run("Valid request", func(t *testing.T) {
		provider.activity = models.WalletActivity{
			Address: "0x123",
			Chain:   "ethereum",
		}
		provider.err = nil

		reqBody, _ := json.Marshal(models.AnalyzeWalletRequest{
			Address: "0x123",
			Chain:   "ethereum",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze-wallet", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.AnalyzeWallet()(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", w.Code)
		}

		var resp models.AnalyzeWalletResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Address != "0x123" {
			t.Errorf("expected address 0x123, got %s", resp.Address)
		}
	})

	t.Run("Empty address", func(t *testing.T) {
		reqBody, _ := json.Marshal(models.AnalyzeWalletRequest{
			Address: "",
			Chain:   "ethereum",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze-wallet", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.AnalyzeWallet()(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest, got %v", w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze-wallet", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.AnalyzeWallet()(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest, got %v", w.Code)
		}
	})

	t.Run("Service error", func(t *testing.T) {
		provider.err = errors.New("some error")

		reqBody, _ := json.Marshal(models.AnalyzeWalletRequest{
			Address: "0x123",
			Chain:   "ethereum",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/analyze-wallet", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler.AnalyzeWallet()(w, req)

		if w.Code != http.StatusBadGateway {
			t.Errorf("expected status BadGateway, got %v", w.Code)
		}
	})
}
