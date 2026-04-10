package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Alaghal/ai-wallet-analyzer/internal/models"
	"github.com/Alaghal/ai-wallet-analyzer/internal/service"
)

type WalletHandler struct {
	analyzer *service.AnalyzerService
}

func NewWalletHandler(analyzer *service.AnalyzerService) *WalletHandler {
	return &WalletHandler{
		analyzer: analyzer,
	}
}

func (h *WalletHandler) AnalyzeWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AnalyzeWalletRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid JSON request body",
			})
			return
		}

		req.Address = strings.TrimSpace(req.Address)
		req.Chain = strings.TrimSpace(req.Chain)

		if req.Address == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "address is required",
			})
			return
		}

		result, err := h.analyzer.Analyze(r.Context(), req.Address, req.Chain)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}
