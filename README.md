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