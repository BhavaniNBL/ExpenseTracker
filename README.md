# ExpenseTracker
Here is your project structure in **clean list format**:

### 📁 Project Folder Structure – `GoExpenseTracker/`

* `cmd/` – Main server entry point
* `config/` – Application configuration and environment loading
* `docs/` – Swagger auto-generated documentation
* `handler/` – HTTP handlers using Gin
* `middleware/` – JWT, logging, and authentication middleware
* `models/` – Struct definitions for expenses and users
* `repository/` – Database access layer (uses GORM)
* `services/` – Core business logic layer
* `go.mod / go.sum` – Go module files for dependency management
* `main.go` – Entry point to bootstrap and start the server



# Using Docker
docker-compose up --build


### Prerequisites

- Go 1.19+
- PostgreSQL
- Docker & Docker Compose

## Setup

### 1. Clone and Install

```bash
git clone https://github.com/yourusername/go-expense-tracker.git
cd go-expense-tracker
go mod tidy

# Run Application
go run main.go

The server runs at http://localhost:8085


# Swagger docs:
http://localhost:8085/swagger/index.html

# API Endpoints
POST    http://localhost:8085/api/v1/auth/login
POST    http://localhost:8085/api/v1/expenses
GET     http://localhost:8085/api/v1/expenses/{expense_id}
PUT     http://localhost:8085/api/v1/expenses/{expense_id}
DELETE  http://localhost:8085/api/v1/expenses/{expense_id}
GET     http://localhost:8085/api/v1/expenses?category=Food&limit=10&offset=0
GET     http://localhost:8085/api/v1/expenses/summary

 
# Unit Testing 

# Run Unit Tests with Coverage
go test ./...  -coverprofile=coverage.out  


# View Coverage in Browser
go tool cover -html="coverage.out"    
