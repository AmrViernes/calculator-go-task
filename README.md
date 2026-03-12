# Pack Calculator API

A flexible REST API for calculating optimal pack sizes for customer orders.

## Features

- Flexible pack size configuration via API or environment variables
- Dynamic programming algorithm for optimal pack combinations
- RESTful JSON API
- Web UI for testing
- Docker deployment ready

## Quick Start

```bash
docker build -t pack-calculator .
docker run -p 8080:8080 pack-calculator
```

Or use docker-compose:
```bash
docker-compose up
```

## API Endpoints

**Calculate Packs**
```http
POST /api/calculate
{"orderQuantity": 501}
```
```json
{"packs": [{"size": 250, "count": 1}, {"size": 500, "count": 1}]}
```

**Get Pack Sizes**
```http
GET /api/packsizes
```

**Update Pack Sizes**
```http
PUT /api/packsizes
{"packSizes": [100, 200, 500, 1000]}
```

## Project Structure

```
pack-calculator/
├── calculator/       # Core algorithm
├── api/             # HTTP handlers
├── ui/              # Web interface
├── main.go          # Entry point
├── Dockerfile
├── docker-compose.yml
└── .render.yaml
```