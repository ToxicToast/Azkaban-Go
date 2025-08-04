# Azkaban 🧙‍♂️

Welcome to the **Azkaban** – a modern, event-driven microservice architecture built with Go, designed for scalability, observability, and maintainability.

## 📦 Project Overview

This monorepo contains all Azkaban microservices, shared libraries, and infrastructure components. Each service is written in **Go** and follows a clear separation of concerns, organized as follows:

```
apps/
  gateway/          # API Gateway "Dementor" – REST entrypoint + gRPC routing + AuthZ
  foodfolio/        # Food management service – ingredients, recipes, shopping list
  warcraft/         # WoW character tracking service (name TBD)
libs/
  shared/           # Common utilities, error handling, etc.
  proto/            # Protobuf definitions for gRPC
  events/           # Kafka event schemas and enums
```

## 🚀 Features

- 🔗 **gRPC & REST**: REST entry via gateway, internal communication via gRPC
- 🔐 **Authentication via Authentik** with JWT group validation in the gateway
- 💬 **Kafka-based event system** for loosely coupled communication
- 🧠 **Redis caching** for faster response times
- 📊 **Observability**: Tracing & logging with OpenTelemetry and SigNoz
- ☠️ **Circuit Breaker** for resilient service-to-service calls

## 🛠 Tech Stack

| Area               | Technology                      |
|--------------------|----------------------------------|
| Language           | Golang                          |
| Communication      | gRPC, Kafka                     |
| Auth               | Authentik (OAuth2 / JWT)        |
| Monitoring         | OpenTelemetry, SigNoz           |
| Gateway            | Custom Go service (Dementor)    |
| Containerization   | Docker, Kubernetes              |
| Caching            | Redis                           |
| Events             | Kafka                           |
| UI (Dashboards)    | React, Redux Toolkit            |

## 📂 Structure

```text
azkaban/
├── apps/
│   ├── gateway/
│   ├── foodfolio/
│   └── warcraft/
├── libs/
│   ├── proto/
│   ├── events/
│   └── shared/
└── README.md
```

## 🧪 Local Development

```bash
# Start the gateway
cd apps/gateway
go run main.go

# Start Foodfolio
cd apps/foodfolio
go run main.go
```

> Note: Local `.env` files are stored within each service directory.

## 🧭 Roadmap

| Sprint | Topic                              |
|--------|------------------------------------|
| 1      | Gateway + Auth + Tracing           |
| 2      | Service registry architecture      |
| 3      | Kafka-based event system           |
| 4      | Event-to-REST adapter              |
| 5      | Event dashboards with React        |

## 👤 Maintainer

**ToxicToast**  
Lead Consultant · Frankfurt, Germany  
> Questions? Feature requests? → [Open an issue](https://github.com/ToxicToast/Azkaban-Go/issues)

## 📜 License

This project is **closed source**. Unauthorized distribution, modification, or commercial use is strictly prohibited.
