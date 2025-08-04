# Proto Definitions

This library contains all Protobuf definitions used in the Azkaban ecosystem.

## Structure

- `*.proto` files for all services
- Shared message formats and enums
- gRPC service definitions

## Usage

```bash
protoc --go_out=. --go-grpc_out=. your_file.proto
```

Use a central generation approach for consistency across services.
