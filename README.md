# Lead Search Super

An intelligent lead enrichment system that combines web scraping, AI-powered data extraction, and microservices architecture to automatically gather and process business information.

## 🚀 Features

- **Automated Lead Enrichment**: Searches and enriches business leads with data from multiple sources
- **AI-Powered Extraction**: Uses Ollama LLM to intelligently extract structured data (CNPJ, company names, contacts) from unstructured web content
- **Microservices Architecture**: Distributed system with specialized services for different tasks
- **Real-time Processing**: RabbitMQ-based message queue for asynchronous processing
- **Multiple Data Sources**: Integrates with Tavily, Serper, and CNPJ.BIZ APIs
- **Data Lake Storage**: Elasticsearch for storing and querying enriched lead data
- **Web Interface**: React-based frontend for searching and viewing leads

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│  Front-Web  │────▶│     API     │────▶│  PostgreSQL  │
└─────────────┘     └──────┬──────┘     └──────────────┘
                           │
                    ┌──────▼──────┐
                    │  RabbitMQ   │
                    │   (Fanout)   │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
┌───────▼──────┐  ┌────────▼────────┐  ┌─────▼──────┐
│Data Collector│  │    Forwarder    │  │  DataLake  │
│              │  │   (Ollama AI)   │  │            │
└───────┬──────┘  └────────┬────────┘  └─────┬──────┘
        │                  │                  │
┌───────▼──────┐  ┌────────▼────────┐  ┌─────▼──────┐
│ Tavily/Serper│  │     Ollama      │  │Elasticsearch│
└──────────────┘  └─────────────────┘  └────────────┘
```

## 🛠️ Services

### API Service (Go)
- REST API for lead management
- PostgreSQL integration
- RabbitMQ producer for lead processing

### Data Collector (Go)
- Enriches leads with external data sources
- Integrates with Tavily, Serper, and CNPJ.BIZ APIs
- Publishes enriched data to RabbitMQ fanout exchange

### Forwarder (Go)
- Consumes enriched lead data
- Uses Ollama AI (qwen2.5:14b model) to extract structured information
- Updates lead records with extracted CNPJ, company names, and other details

### DataLake (Go)
- Stores enriched lead data in Elasticsearch
- Provides searchable archive of all processed leads

### Front-Web (React + TypeScript)
- User interface for searching and viewing leads
- Real-time updates of lead processing status
- Google Places integration for initial lead discovery

### Search-Google (Node.js)
- Google Places API integration
- Initial lead discovery service

## 📋 Prerequisites

- Docker and Docker Compose
- Ollama installed locally with qwen2.5:14b model
- API Keys for:
  - Google Places API
  - Tavily API
  - Serper API
  - CNPJ.BIZ API (optional)

## 🚀 Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/lead-search-super.git
cd lead-search-super
```

2. **Set up environment variables**
```bash
# Create .env files in each service directory with required API keys
cp .env.example .env
```

3. **Install Ollama and pull the model**
```bash
ollama pull qwen2.5:14b
```

4. **Start the services**
```bash
docker-compose up -d
```

5. **Access the application**
- Frontend: http://localhost:5173
- API: http://localhost:8085
- RabbitMQ Management: http://localhost:15672 (guest/guest)
- Elasticsearch: http://localhost:9200

## 🔧 Configuration

### Docker Compose Environment Variables

- `RABBITMQ_URL`: RabbitMQ connection string
- `OLHAMA_URL`: Ollama API endpoint (use `http://host.docker.internal:11434/api/chat` for local Ollama)
- `ELASTICSEARCH_URL`: Elasticsearch connection URL
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`: PostgreSQL credentials

### API Keys Configuration

Each service requires specific API keys in their `.env` files:

- **search-google**: `GOOGLE_PLACES_API_KEY`
- **data-collector**: `TAVILY_API_KEY`, `SERPER_API_KEY`
- **forwarder**: `INVERTEXTO_API_TOKEN` (optional)

## 📊 Data Flow

1. User searches for a business in the frontend
2. Search-Google service queries Google Places API
3. API service saves the lead and publishes to RabbitMQ
4. Data Collector enriches the lead with Tavily and Serper data
5. Enriched data is published to fanout exchange
6. Forwarder uses Ollama AI to extract structured data (CNPJ, company details)
7. DataLake stores the complete enriched data in Elasticsearch
8. Frontend displays updated lead information

## 🧪 Testing

### Manual Lead Publishing
```bash
# Use the provided script to manually publish a lead
./publish_lead.sh
```

### Check Processing Status
```bash
# View RabbitMQ queue status
docker exec rabbitmq rabbitmqctl list_queues

# Check forwarder logs
docker logs forwarder --tail 100

# Query processed leads
docker exec db-leads psql -U leads_user -d leads_db -c "SELECT * FROM leads;"
```

## 🐛 Troubleshooting

### Ollama Connection Issues
- Ensure Ollama is running locally: `ollama serve`
- Verify the model is installed: `ollama list`
- Check Docker can reach host: `http://host.docker.internal:11434`

### Message Processing Stuck
- Check RabbitMQ for unacknowledged messages
- Restart the forwarder service: `docker restart forwarder`
- Check logs for errors: `docker logs forwarder`

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 👥 Authors

- Bruno Vieira

## 🙏 Acknowledgments

- Ollama for providing local LLM capabilities
- Tavily and Serper for search APIs
- The Go and React communities for excellent tools and libraries