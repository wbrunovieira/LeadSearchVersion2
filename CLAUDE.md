# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based lead search and enrichment system with the following key components:

### Services
- **api** (Go, port 8085): REST API service with PostgreSQL database integration and RabbitMQ messaging
- **data-collector** (Go, port 8086): Collects data from various sources (CNPJ, Serper, Tavily)
- **search-google** (Go, port 8082): Google Places API integration service
- **forwarder** (Go): Processes messages from RabbitMQ and enriches lead data using Ollama LLM
- **datalake** (Go, port 8087): Stores enriched data in Elasticsearch
- **front-web** (React/TypeScript/Vite, port 5173): Web frontend application

### Infrastructure
- **PostgreSQL**: Primary database for leads (db-leads)
- **RabbitMQ**: Message broker for service communication (ports 5672, 15672)
- **Elasticsearch**: Data lake storage for enriched leads (port 9200)

## Common Development Commands

### Frontend (front-web)
```bash
cd front-web
pnpm install          # Install dependencies
pnpm dev             # Start development server (port 5173)
pnpm build           # Build for production
pnpm lint            # Run ESLint
```

### Backend Services (Go)
```bash
# Run tests for a specific service
cd [service-directory]
go test ./...

# Build a service
go build -o bin/[service-name] main.go

# Run a service locally
go run main.go
```

### Docker Operations
```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d [service-name]

# View logs
docker-compose logs -f [service-name]

# Rebuild and restart a service
docker-compose up -d --build [service-name]

# Stop all services
docker-compose down
```

## Service Communication Flow

1. Frontend (front-web) → API service for lead management
2. API → RabbitMQ → data-collector for enrichment requests
3. data-collector → External APIs (CNPJ, Serper, Tavily) for data collection
4. RabbitMQ → forwarder → Ollama LLM for data enrichment
5. Enriched data → datalake → Elasticsearch for storage

## Key Implementation Details

### API Service
- Database migrations handled automatically on startup
- CORS middleware enabled for frontend communication
- Endpoints: `/save-leads`, `/update-lead-field`, `/health`

### RabbitMQ Queues
- Services communicate through RabbitMQ message queues
- Connection URL: `amqp://guest:guest@rabbitmq:5672/`

### Frontend Stack
- React 19 with TypeScript
- Vite for build tooling
- Tailwind CSS v4 for styling
- ESLint configured for code quality

### Environment Variables
- Services use environment variables for configuration
- Key variables: `PORT`, `RABBITMQ_URL`, `DB_*` (database config), `ELASTICSEARCH_URL`
- Ollama LLM endpoint configured in forwarder: `OLHAMA_URL`