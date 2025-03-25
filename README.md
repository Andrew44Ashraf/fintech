#  Fintech Service

A Go microservice for account and transaction management with PostgreSQL.

---

##  Quick Start

```bash
# Clone the repository
git clone https://github.com/Andrew44Ashraf/fintech-service.git
cd fintech-service

# Run the service using Docker
docker-compose up --build

--- 
## Installation

git clone https://github.com/Andrew44Ashraf/fintech-service.git
cd fintech-service

# Initialize Go modules
go mod init github.com/Andrew44Ashraf/fintech-service
go mod tidy

# Run the application using Docker
docker-compose up --build



Local Development
Start PostgreSQL (either locally or using Docker)

docker-compose up -d db
Run Migrations


go run cmd/migrate/main.go
Start the Application


go run cmd/main.go
📌 Prerequisites
Go 1.21+

Docker & Docker Compose
