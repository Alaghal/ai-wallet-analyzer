package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Alaghal/ai-wallet-analyzer/internal/ai"
	"github.com/Alaghal/ai-wallet-analyzer/internal/config"
	"github.com/Alaghal/ai-wallet-analyzer/internal/handlers"
	appMetrics "github.com/Alaghal/ai-wallet-analyzer/internal/metrics"
	appMiddleware "github.com/Alaghal/ai-wallet-analyzer/internal/middleware"
	"github.com/Alaghal/ai-wallet-analyzer/internal/provider"
	"github.com/Alaghal/ai-wallet-analyzer/internal/service"
)

type Server struct {
	httpServer *http.Server
	cfg        config.Config
}

func New(cfg config.Config) *Server {
	router := newRouter(cfg)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.AppPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		cfg:        cfg,
	}
}

func newRouter(cfg config.Config) http.Handler {
	router := chi.NewRouter()

	metricsRegistry := appMetrics.MustNew()

	activityProvider, providerName := buildProvider(cfg)
	llmClient, llmName := buildLLMClient(cfg)

	instrumentedProvider := provider.NewInstrumentedProvider(providerName, activityProvider, metricsRegistry)
	instrumentedLLM := ai.NewInstrumentedClient(llmName, llmClient, metricsRegistry)

	analyzerService := service.NewAnalyzerService(instrumentedProvider, instrumentedLLM)
	walletHandler := handlers.NewWalletHandler(analyzerService)

	router.Use(chiMiddleware.Timeout(15 * time.Second))
	router.Use(appMiddleware.RequestID)
	router.Use(appMiddleware.Recovery)
	router.Use(appMiddleware.Logging)
	router.Use(appMiddleware.Metrics(metricsRegistry))

	router.Get("/health", handlers.Health())
	router.Handle("/metrics", promhttp.Handler())

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/analyze-wallet", walletHandler.AnalyzeWallet())
		})
	})

	return router
}

func buildProvider(cfg config.Config) (provider.WalletActivityProvider, string) {
	switch cfg.ProviderType {
	case "etherscan":
		log.Printf("using wallet activity provider=etherscan")
		return provider.NewEtherscanWalletActivityProvider(
			cfg.EtherscanAPIURL,
			cfg.EtherscanAPIKey,
			cfg.HTTPTimeout,
		), "etherscan"
	default:
		log.Printf("using wallet activity provider=mock")
		return provider.NewMockWalletActivityProvider(), "mock"
	}
}

func buildLLMClient(cfg config.Config) (ai.Client, string) {
	switch cfg.LLMProviderType {
	case "openai":
		log.Printf("using llm provider=openai-compatible model=%s", cfg.OpenAIModel)
		return ai.NewOpenAIClient(
			cfg.OpenAIAPIURL,
			cfg.OpenAIAPIKey,
			cfg.OpenAIModel,
			cfg.HTTPTimeout,
		), "openai"
	default:
		log.Printf("using llm provider=mock")
		return ai.NewMockClient(), "mock"
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("starting ai-wallet-analyzer on port %d (env=%s)", s.cfg.AppPort, s.cfg.AppEnv)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("server stopped gracefully")
		return nil

	case err := <-errCh:
		return fmt.Errorf("http server failed: %w", err)
	}
}
