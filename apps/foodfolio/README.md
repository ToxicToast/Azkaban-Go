# Foodfolio Service

This microservice handles food inventory and consumption tracking, inspired by tools like Grocy.

## Responsibilities

- Manage ingredients, shopping lists, and meals
- Publish food-related events to Kafka
- Persist data in PostgreSQL
- Expose gRPC API (and REST via gateway)

## Stack

- Golang 1.24
- PostgreSQL
- gRPC
- Redis (cache)
- Kafka (event streaming)