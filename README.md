# csp-report-collector

[![Docker Image](https://img.shields.io/badge/docker-ghcr.io%2Fmeysam81%2Fcsp--report--collector-blue)](https://github.com/meysam81/csp-report-collector/pkgs/container/csp-report-collector)
[![Go Report Card](https://goreportcard.com/badge/github.com/meysam81/csp-report-collector)](https://goreportcard.com/report/github.com/meysam81/csp-report-collector)
[![License](https://img.shields.io/github/license/meysam81/csp-report-collector)](LICENSE)

Lightweight service to collect and persist Content Security Policy (CSP) violation reports in Redis for audit and investigation.

## Quick Start

```bash
docker run --rm -dp 8080:8080 \
  -e REDIS_HOST=your-redis-host \
  ghcr.io/meysam81/csp-report-collector
```

## Configure Your CSP Header

Point your CSP reporting to the collector:

```shell
Content-Security-Policy: default-src 'self'; report-uri https://your-domain.com:8080/
```

Or with the Reporting API:

```shell
Content-Security-Policy: default-src 'self'; report-to csp-endpoint
Report-To: {"group":"csp-endpoint","max_age":86400,"endpoints":[{"url":"https://your-domain.com:8080/"}]}
```

## Configuration

All configuration follows the 12-factor app methodology via environment variables:

```shell
Content-Security-Policy: default-src 'self'; report-to csp-endpoint
Report-To: {"group":"csp-endpoint","max_age":86400,"endpoints":[{"url":"https://your-domain.com:8080/"}]}
```

## Configuration

All configuration follows the 12-factor app methodology via environment variables:

| Variable             | Default     | Description                           |
| -------------------- | ----------- | ------------------------------------- |
| `PORT`               | `8080`      | HTTP server port                      |
| `LOG_LEVEL`          | `info`      | Log verbosity (debug/info/warn/error) |
| `REDIS_HOST`         | `localhost` | Redis server hostname (**required**)  |
| `REDIS_PORT`         | `6379`      | Redis server port                     |
| `REDIS_DB`           | `0`         | Redis database number                 |
| `REDIS_PASSWORD`     | -           | Redis authentication password         |
| `REDIS_SSL__ENABLED` | `false`     | Enable TLS connection to Redis        |
| `RATELIMIT_MAX`      | `20`        | Max requests per IP                   |
| `RATELIMIT_REFILL`   | `2.0`       | Token refill rate per second          |

## Data Storage

CSP reports are stored in Redis with:

- **Key**: Unix timestamp of receipt
- **Value**: Full JSON report
- **TTL**: Indefinite (configure Redis eviction policy as needed)

## Example Report Format

The collector accepts standard CSP violation reports:

```json
{
  "age": 53531,
  "body": {
    "blockedURL": "inline",
    "disposition": "enforce",
    "documentURL": "https://example.com/page",
    "effectiveDirective": "script-src-elem",
    "originalPolicy": "default-src 'self'",
    "statusCode": 200
  },
  "type": "csp-violation",
  "url": "https://example.com/page",
  "user_agent": "Mozilla/5.0..."
}
```

## Docker Compose Example

```yaml
version: "3.8"
services:
  csp-collector:
    image: ghcr.io/meysam81/csp-report-collector
    ports:
      - "8080:8080"
    environment:
      REDIS_HOST: redis
      RATELIMIT_MAX: 50
    depends_on:
      - redis

  redis:
    image: redis:alpine
    volumes:
      - redis-data:/data

volumes:
  redis-data:
```

## Rate Limiting

Built-in rate limiting per IP address returns:

- `X-RateLimit-Total` header with limit
- `X-RateLimit-Remaining` header with remaining requests
- `429 Too Many Requests` when exceeded

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
