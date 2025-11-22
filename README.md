# Project_A

A REST API service that provides tide timing information for Singapore by scraping data from the National Environment Agency (NEA) website.

## Features

- **Tide Timings API**: Retrieve comprehensive tide information including high/low tide times, heights, and levels
- **Health Check**: Built-in health endpoint to monitor service status and uptime
- **Docker Support**: Containerized deployment with multi-stage Docker build
- **Clean Architecture**: Well-structured Go application with separate layers for handlers, services, and data transfer objects

## API Endpoints

### Health Check
- `GET /check-health` - Returns service health status and uptime

### Tide Timings
- `GET /tide-timings` - Returns tide timing data for available months

### Response Format

The tide timings endpoint returns data in the following JSON structure:

```json
{
  "data": [
    {
      "month": "January 2025",
      "days": [
        {
          "day": 1,
          "observations": [
            {
              "time": "02:15",
              "height": 2.8,
              "level": "High Tide"
            },
            {
              "time": "08:45",
              "height": 0.9,
              "level": "Low Tide"
            }
          ]
        }
      ]
    }
  ]
}
```

## Prerequisites

- Go 1.25.4 or later
- Docker (optional, for containerized deployment)

## Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd project-a
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Environment Configuration**
   Create a `.env` file in the project root:
   ```env
   PORT=8080
   SECRET_JWT=your-secret-key-here
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

5. **Test the API**
   ```bash
   curl http://localhost:8080/check-health
   curl http://localhost:8080/tide-timings
   ```

## Docker Deployment

### Using Docker Compose

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

2. **Access the API**
   - Health check: `http://localhost:8080/check-health`
   - Tide timings: `http://localhost:8080/tide-timings`

### Using Docker directly

1. **Build the image**
   ```bash
   docker build -t project-a .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 -e PORT=8080 project-a
   ```

## Project Structure

```
project-a/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── dto/                    # Data transfer objects
│   ├── handler/                # HTTP request handlers
│   ├── service/                # Business logic layer
│   └── middleware/             # HTTP middleware
├── pkg/                        # Shared packages
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # Docker Compose configuration
├── go.mod                      # Go module dependencies
└── go.sum                      # Dependency checksums
```

## Architecture

The application follows a clean architecture pattern with three main layers:

- **Handler Layer** (`internal/handler/`): HTTP request/response handling using Gin framework
- **Service Layer** (`internal/service/`): Business logic and external API interactions
- **Data Layer** (`internal/dto/`): Data structures for API responses

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing and scraping
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol (for future enhancements)

## Data Source

Tide timing data is sourced from the Singapore National Environment Agency (NEA) website:
https://www.nea.gov.sg/corporate-functions/weather/tide-timings

The service scrapes and parses HTML content to extract tide information including:
- Monthly tide schedules
- Daily tide observations
- Tide heights and levels (High/Low)

## Error Handling

The API provides appropriate HTTP status codes:
- `200 OK` - Successful response
- `502 Bad Gateway` - External data source unavailable or parsing errors

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]

