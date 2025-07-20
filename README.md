# Search-Proxy

A proxy service for Brave Search API that helps bypass rate limits for Open WebUI search requests.

## Overview

Search-Proxy is a lightweight Go application that acts as an intermediary between your application and the Brave Search API. It provides rate limiting capabilities to manage API request traffic and prevent hitting Brave's rate limits.

## Features

- Proxies search requests to Brave Search API
- Built-in token-based rate limiting
- Configurable concurrency and cooldown settings
- Simple JSON API for search queries
- Modular architecture for potential expansion to other search providers

## Requirements

- Go 1.23 or higher
- Brave Search API key (obtain from [Brave Search API](https://brave.com/search/api/))

## Installation

### Clone the repository

```bash
git clone https://github.com/lokman928/search-proxy.git
cd search-proxy
```

### Build the application

```bash
make build
```

This will create a binary in the `bin` directory.

## Configuration

Copy the example configuration file and modify it with your settings:

```bash
cp config.example.toml config.toml
```

Edit `config.toml` with your preferred settings:

```toml
[server]
port = 8080  # Port the proxy server will listen on

[brave]
base_url = "https://api.search.brave.com"  # Brave Search API base URL
api_key = "<your api key here>"  # Your Brave Search API key

# Rate limiting configuration
rate_limit.enable = true  # Enable/disable rate limiting
rate_limit.cooldown_time = 1  # Cooldown time in milliseconds
rate_limit.max_concurrency = 1  # Maximum concurrent requests
```

## Usage

### Running the server

```bash
make run
```

Or directly run the binary:

```bash
./bin/proxy
```

### API Endpoints

#### Search

```
POST /brave/search
```

Request body:

```json
{
  "query": "your search query",
  "count": 10  # Number of results to return
}
```

Response:

```json
[
  {
    "link": "https://example.com",
    "title": "Example Website",
    "snippet": "This is an example search result snippet."
  },
  ...
]
```

## Development

### Running tests

```bash
make unittest
```

### Project Structure

- `cmd/proxy/` - Application entry point
- `internal/` - Internal packages
  - `common/` - Shared models and utilities
  - `config/` - Configuration handling
  - `module/` - Modular components (currently Brave search)
  - `proxy/` - Main proxy server implementation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
