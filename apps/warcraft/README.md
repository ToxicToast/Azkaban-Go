# Warcraft Service

This service handles World of Warcraft character tracking and event collection.

## Features

- Sync characters from Blizzard API
- Track item level, raid progress, guilds
- Publish game events to Kafka

## Stack

- Golang 1.24
- PostgreSQL
- gRPC
- Redis (cache)
- Kafka (event streaming)