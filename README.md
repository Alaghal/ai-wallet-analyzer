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