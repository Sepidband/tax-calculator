# Tax Calculator

A Canadian income tax calculator with REST API and web interface.

## Features

-  Calculate total income tax using marginal tax rates
-  Breakdown by tax brackets
-  Effective tax rate calculation
-  REST API with OpenAPI documentation
-  Simple web interface
-  Automatic retry logic for API failures
-  Comprehensive logging

## Quick Start

### Prerequisites

- Go 1.19+
- Docker (for tax API service)

### Running the Application

1. **Start the tax data API:**
```bash
   docker pull ptsdocker16/interview-test-server
   docker run --init -p 5001:5001 -it ptsdocker16/interview-test-server
```

2. **Start the backend API:**
```bash
    git clone https://github.com/Sepidband/tax-calculator.git
    cd tax-calculator
    go mod download
    go run cmd/server/main.go
```
3. **Start the frontend (separate terminal):**
```bash 
    cd frontend
   python -m http.server 3000
```
4. **Access the application:**

    - Web UI: http://localhost:3000
    - Health Check: http://localhost:8080/health

## API Usage
### Calculate Tax
#### POST `/api/v1/calculate-tax`
```bash
{
    "salary": 50000,
    "year": 2022
}
```
#### Response:
```bash
{
  "total_tax": 7500.00,
  "effective_rate": 0.15,
  "bracket_breakdown": [
    {
      "bracket": {
        "min": 0,
        "max": 50197,
        "rate": 0.15
      },
      "tax_owed": 7500.00
    }
  ]
}
```
### Example with cURL
```bash
curl -X POST http://localhost:8080/api/v1/calculate-tax \
  -H "Content-Type: application/json" \
  -d '{"salary": 50000, "year": 2022}'
```
### Architecture
```bash
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │    │   Tax Data API  │
│   (JavaScript)  │───▶│   (Go/Gin)      │───▶│   (Docker)      │
│   Port 3000     │    │   Port 8080     │    │   Port 5001     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```


### Project Structure
```bash
tax-calculator/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── api/handlers.go         # HTTP handlers
│   ├── calculator/calculator.go # Tax calculation logic
│   └── client/tax_api.go       # External API client
├── pkg/models/tax.go           # Data models
├── frontend/                   # Web interface
│   ├── index.html
│   ├── script.js
│   └── style.css
├── tests/                      # Test files
├── docs/                       # Generated API docs
└── README.md
```

### Testing
```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestCalculateTax ./tests
```

### Configuration

#### Environment Variables

| Variable | Type | Default | Required | Description | Example |
|----------|------|---------|----------|-------------|---------|
| `TAX_API_URL` | string | `http://localhost:5001` | No | Base URL for the external tax data API service | `http://tax-api:5001` |
| `PORT` | string | `8080` | No | Port number for the HTTP server to listen on | `3000` |
| `LOG_LEVEL` | string | `info` | No | Logging level (trace, debug, info, warn, error, fatal, panic) | `debug` |

### Tax Calculation Logic
This application implements Canadian marginal tax rates:

1. Marginal Taxation: Each bracket is taxed at its respective rate
2. Progressive System: Higher income portions are taxed at higher rates
3. Effective Rate: Total tax divided by total income

#### Example Calculation for $100,000 (2022 rates):
- Bracket 1: $0 - $50,197 at 15% = $7,529.55
- Bracket 2: $50,197 - $100,000 at 20.5% = $10,209.62
- Total Tax: $17,739.17
- Effective Rate: 17.74%

### Error Handling
The application handles several error scenarios:

- Invalid year: Only 2019-2022 supported
- API failures: Automatic retry with exponential backoff
- Invalid input: Proper validation and user-friendly messages
- Network issues: Graceful degradation

### Production Considerations

- Logging: Structured logging with different levels
- Monitoring: Health check endpoint for load balancers
- Security: CORS configured, input validation
- Performance: Connection pooling, request timeouts
- Scalability: Stateless design, horizontal scaling ready