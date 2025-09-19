# Tax Calculator

A Canadian income tax calculator with REST API and web interface featuring in-memory caching and interactive API documentation.

## Features

-  Calculate total income tax using marginal tax rates
-  Breakdown by tax brackets
-  Effective tax rate calculation
-  REST API with OpenAPI documentation
-  Simple web interface
-  Automatic retry logic for API failures
-  Comprehensive logging
-  In-memory caching for improved performance
-  Interactive Swagger/OpenAPI documentation

## Quick Start

### Prerequisites

- Go 1.21+
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
    - API Documentation: http://localhost:8080/swagger/index.html
    - Health Check: http://localhost:8080/api/v1/health

## API Documentation
### Interactive Documentation
The API includes comprehensive Swagger/OpenAPI documentation accessible at:
```bash
http://localhost:8080/swagger/index.html
```
### Features:

- Interactive API explorer - Test endpoints directly from the browser
- Request/response schemas - Complete data model documentation
- Example requests - Sample data for easy testing
- Error responses - Detailed error handling documentation

### Generating Documentation
Documentation is automatically generated from code annotations. To regenerate:
```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/server/main.go -d ./,./internal/api,./pkg/models

# Documentation files are created in docs/
```

## Caching
### Overview
The application implements intelligent caching to improve performance and reliability:

- Tax brackets caching - Reduces external API calls
- Automatic cache invalidation - Ensures data freshness
- Fallback handling - Graceful degradation during API failures

### Cache Configuration

| Variable | Type | Default | Required | Description | Example |
|----------|------|---------|----------|-------------|---------|
| `CACHE_ENABLED` | bool | `true` | No | Enable or disable in-memory caching system | `false` |
| `CACHE_TTL` | duration | `24h` | No | Time-to-live for cached tax bracket data | `12h`, `2d`, `1h30m` |
| `CACHE_MAX_SIZE` | int | `100` | No | Maximum number of cached entries (tax years) | `50`, `200` |
| `CACHE_CLEANUP_INTERVAL` | duration | `1h` | No | How often to run cache cleanup for expired entries | `30m`, `2h` |
| `CACHE_WARMUP_ENABLED` | bool | `true` | No | Pre-populate cache with common tax years on startup | `false` |
| `CACHE_WARMUP_YEARS` | string | `2019,2020,2021,2022` | No | Comma-separated list of years to pre-cache | `2021,2022,2023` |
| `CACHE_METRICS_ENABLED` | bool | `true` | No | Enable cache performance metrics collection | `false` |
| `CACHE_LOG_HITS` | bool | `false` | No | Log cache hits and misses for debugging | `true` |


### Cache Behavior 
```bash
Request Flow with Caching:
┌─────────────┐    ┌──────────┐    ┌─────────────┐    ┌─────────────┐
│   Request   │───▶│  Cache   │───▶│ External    │───▶│  Response   │
│             │    │  Check   │    │ API Call    │    │             │
└─────────────┘    └──────────┘    └─────────────┘    └─────────────┘
                        │                                     │
                        ▼                                     │
                   Cache Hit? ──────────────────────────────────┘
                   Return Cached Data
```


### Cache Strategy:

- First request: Fetches from external API, stores in cache
- Subsequent requests: Returns cached data (if valid)
- Cache miss/expiry: Automatically refetches and updates cache
- API failures: Uses cached data as fallback (if available)

### Performance Benefits

- Reduced latency: ~95% faster response times for cached requests
- Lower external API load: Significant reduction in upstream calls
- Improved reliability: Continues working during external API outages
- Cost efficiency: Fewer external API requests

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

## Health Check
### GET `/api/v1/health`
#### Response:
```bash
{
  "status": "healthy",
  "service": "tax-calculator",
  "timestamp": 1695234567
}
```

### Example with cURL
```bash
curl -X POST http://localhost:8080/api/v1/calculate-tax \
  -H "Content-Type: application/json" \
  -d '{"salary": 50000, "year": 2022}'
```
## Architecture
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
├── docs/                       # Generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml   └── style.css
├── tests/                      # Test files
└── README.md
```

## Testing
```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestCalculateTax ./tests

# Test with cache scenarios
go test -run TestCache ./internal/client
```

## Configuration

#### Environment Variables

| Variable | Type | Default | Required | Description | Example |
|----------|------|---------|----------|-------------|---------|
| `TAX_API_URL` | string | `http://localhost:5001` | No | Base URL for the external tax data API service | `http://tax-api:5001` |
| `PORT` | string | `8080` | No | Port number for the HTTP server to listen on | `3000` |
| `LOG_LEVEL` | string | `info` | No | Logging level (trace, debug, info, warn, error, fatal, panic) | `debug` |

## Tax Calculation Logic
This application implements Canadian marginal tax rates:

1. Marginal Taxation: Each bracket is taxed at its respective rate
2. Progressive System: Higher income portions are taxed at higher rates
3. Effective Rate: Total tax divided by total income

#### Example Calculation for $100,000 (2022 rates):
- Bracket 1: $0 - $50,197 at 15% = $7,529.55
- Bracket 2: $50,197 - $100,000 at 20.5% = $10,209.62
- Total Tax: $17,739.17
- Effective Rate: 17.74%

## Error Handling
The application handles several error scenarios:

- Invalid year: Only 2019-2022 supported
- API failures: Automatic retry with exponential backoff
- Invalid input: Proper validation and user-friendly messages
- Network issues: Graceful degradation
- Cache misses: Transparent fallback to external API

## Performance Optimizations
### Caching Strategy

- Tax brackets cached for 24 hours (configurable)
- Automatic cache warming on first request per year
- Memory-efficient storage with TTL-based expiration
- Cache statistics available via health endpoint

### Production Considerations

- Logging: Structured logging with different levels
- Monitoring: Health check endpoint for load balancers
- Security: CORS configured, input validation
- Performance: Connection pooling, request timeouts
- Scalability: Stateless design, horizontal scaling ready
- Documentation: Interactive API docs for easy integration