# ai-wallet-analyzer

AI-powered Web3 wallet analysis service written in Go.

## Features

- REST API for wallet analysis
- Mock wallet intelligence response
- Environment-based configuration
- Graceful shutdown
- Clean project structure

## Endpoints

### Health
`GET /health`

### Analyze wallet
`POST /api/v1/analyze-wallet`

Example request:

```json
{
  "address": "0x1234567890abcdef",
  "chain": "ethereum"
}
```
## AI Summary Generation

The service supports pluggable LLM providers for wallet activity summarization.

### Supported LLM providers

- `mock` - local mock summary generator for development
- `openai` - OpenAI-compatible chat completion API

### Configuration

- `LLM_PROVIDER_TYPE` - provider type (`mock` by default)
- `OPENAI_API_URL` - OpenAI-compatible API URL
- `OPENAI_API_KEY` - provider API key
- `OPENAI_MODEL` - model name

### Example

```bash
export LLM_PROVIDER_TYPE=openai
export OPENAI_API_KEY=your_api_key
export OPENAI_MODEL=gpt-4o-mini
go run ./cmd/server
```
## Observability

The service includes built-in observability features:

- request correlation via `X-Request-Id`
- access logging with method, path, status, and duration
- panic recovery middleware
- Prometheus-compatible metrics via `GET /metrics`

### Collected metrics

- HTTP request count
- HTTP request duration
- in-flight HTTP requests
- wallet activity provider request count and latency
- LLM request count and latency

### Example

```bash
curl http://localhost:8080/metrics
```

## Docker Compose

The project includes a local monitoring stack with:

- `ai-wallet-analyzer` service
- `Prometheus` for metrics scraping
- `Grafana` for visualization

### Run locally

```bash
docker compose up --build
```
## Monitoring

The service exposes Prometheus-compatible metrics at:

```text
GET /metrics
```