# Gateway (Dementor)

This is the main API Gateway for the Azkaban microservice architecture. It handles:

- JWT-based authentication (via Authentik)
- Request routing (REST to gRPC)
- Request validation & RBAC
- Observability (logging, tracing, metrics)
- Circuit breaking for unavailable services

## Stack

- Golang 1.24
- Gin
- gRPC + gRPC-Gateway
- Kafka (event forwarding)
- Redis (auth/session caching)
- OpenTelemetry