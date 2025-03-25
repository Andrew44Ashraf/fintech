# Fintech Service

A Go microservice for account and transaction management with PostgreSQL.

## Quick Start

```bash
# Clone and run
git clone https://github.com/Andrew44Ashraf/fintech-service.git
cd fintech-service
docker-compose up --build


Installation
git clone https://github.com/Andrew44Ashraf/fintech-service.git
cd fintech-service
go mod init github.com/Andrew44Ashraf/fintech-service
go mod tidy
docker-compose up --build

Local Development
Start PostgreSQL: locally or using docker.
go run cmd/migrate/main.go
go run cmd/main.go


Prerequisites
Go 1.21+

Docker & Docker Compose
 Installation

# Initialize Go modules
go mod init github.com/Andrew44Ashraf/fintech-service
go mod tidy

# Start services
docker-compose up --build


🛠️ Local Development
# Start just PostgreSQL
docker-compose up -d db

# Run migrations
go run cmd/migrate/main.go

# Start application
go run cmd/main.go
