# 🤖 AI Wallet Analyzer

AI-powered Web3 wallet analysis service written in Go. This service provides insights into wallet activity, risk scores, and transaction history summaries using pluggable LLM providers.

## 🚀 Features

- ⚡ **REST API**: Clean endpoints for wallet analysis and health checks.
- 🔌 **Pluggable Providers**: Supports multiple data sources (Etherscan, Mock) and LLM providers (OpenAI, Mock).
- 📊 **Advanced Observability**: Built-in Prometheus metrics, structured logging with Request ID correlation, and panic recovery.
- 🐳 **Containerized**: Ready-to-use Docker and Docker Compose setup with a full monitoring stack (Prometheus & Grafana).
- 🔌 **Graceful Shutdown**: Handles OS signals for clean resource cleanup.

## 🛠 Tech Stack

- **Language**: Go 1.25+
- **External APIs**: Etherscan API (for real data), OpenAI Chat Completions API.
- **Metrics**: Prometheus.
- **Visualization**: Grafana.

## 🏁 Getting Started

### 📋 Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (optional, for monitoring stack)

### 💻 Local Development

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/ai-wallet-analyzer.git
    cd ai-wallet-analyzer
    ```

2.  **Set environment variables** (optional, uses defaults for mock mode):
    ```bash
    export PROVIDER_TYPE=etherscan
    export ETHERSCAN_API_KEY=your_etherscan_key
    export LLM_PROVIDER_TYPE=openai
    export OPENAI_API_KEY=your_openai_key
    ```

3.  **Run the service**:
    ```bash
    go run ./cmd/server
    ```

### 🚢 Docker Compose

Launch the analyzer along with Prometheus and Grafana:

```bash
docker compose up --build
```

- **Service**: `http://localhost:8080`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (Default login: `admin/admin`)

## ⚙️ Configuration

Configuration is managed via environment variables:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `APP_ENV` | Application environment (`local`, `prod`) | `local` |
| `APP_PORT` | Port to listen on | `8080` |
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) | `info` |
| `PROVIDER_TYPE` | Wallet data provider (`mock`, `etherscan`) | `mock` |
| `ETHERSCAN_API_URL` | Etherscan API base URL | `https://api.etherscan.io/api` |
| `ETHERSCAN_API_KEY` | Your Etherscan API key | (empty) |
| `HTTP_TIMEOUT_SECONDS` | Timeout for external API requests | `5` |
| `LLM_PROVIDER_TYPE` | LLM provider type (`mock`, `openai`) | `mock` |
| `OPENAI_API_URL` | OpenAI-compatible API URL | `https://api.openai.com/v1/chat/completions` |
| `OPENAI_API_KEY` | OpenAI API key | (empty) |
| `OPENAI_MODEL` | AI model to use | `gpt-4o-mini` |

## 🛣 API Endpoints

### 1. 🔍 Analyze Wallet
`POST /api/v1/analyze-wallet`

Analyzes a specific wallet address on a given blockchain.

**Request Body:**
```json
{
  "address": "0x1234567890abcdef1234567890abcdef12345678",
  "chain": "ethereum"
}
```

**Response Body:**
```json
{
  "address": "0x1234567890abcdef1234567890abcdef12345678",
  "chain": "ethereum",
  "riskScore": 15,
  "activityLevel": "High",
  "transactionCount": 150,
  "uniqueInteractions": 42,
  "tokens": ["ETH", "USDT", "LINK"],
  "summary": "This wallet shows high frequency of swaps and stablecoin transfers..."
}
```

### 2. 🏥 Health Check
`GET /health`

Returns service status.

### 3. 📈 Metrics
`GET /metrics`

Exposes Prometheus-compatible metrics.

## 🔭 Observability

- **🆔 Request ID**: Every request is assigned a unique `X-Request-Id` (returned in headers).
- **📝 Access Logs**: Standard structured logging for all HTTP requests.
- **📊 Metrics**:
    - `http_requests_total`: Total count of HTTP requests.
    - `http_request_duration_seconds`: Request latency.
    - `provider_requests_total`: Count of external provider requests.
    - `llm_requests_total`: Count of LLM API requests.

## 🏗 Architecture

```mermaid
flowchart LR
    %% Client
    C[Client] -->|POST /analyze-wallet| API[API Layer]

    %% Middleware
    API --> MW[Middleware]
    MW --> RID[Request ID]
    MW --> LOG[Logging]
    MW --> REC[Recovery]
    MW --> MET[Metrics]

    %% Core service
    API --> S[Analyzer Service]

    %% Provider
    S --> P[Wallet Activity Provider]
    P -->|mock| MP[Mock Provider]
    P -->|etherscan| EP[Etherscan API]

    %% Processing
    S --> PROC[Analysis Pipeline]
    PROC --> SCORE[Risk Score]
    PROC --> LEVEL[Activity Level]

    %% LLM
    PROC --> PROMPT[Prompt Builder]
    PROMPT --> LLM[LLM Client]
    LLM -->|mock| ML[Mock LLM]
    LLM -->|openai| OL[OpenAI API]

    %% Response
    LLM --> RESP[Summary]
    RESP --> OUT[JSON Response]
    OUT --> C

    %% Metrics
    MET --> METRICS[/Metrics Endpoint]
    METRICS --> PROM[Prometheus]
    PROM --> GRAF[Grafana]
```