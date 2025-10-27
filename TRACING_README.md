# Distributed Tracing Setup

This project implements distributed tracing using OpenTelemetry and Jaeger.

## Architecture

- **OpenTelemetry**: Standard observability framework for instrumenting applications
- **Jaeger**: Distributed tracing backend for collecting, storing, and visualizing traces

## Components

### 1. Jaeger Service (docker-compose.yml)

Jaeger is configured as an all-in-one instance with:
- **UI Port**: `16686` - Visualize traces in the web interface
- **Collector Port**: `14268` - HTTP endpoint for receiving traces
- **Agent Ports**: `6831/6832` (UDP) - For receiving traces from services

### 2. Blog Service Instrumentation

The blog-service is instrumented with OpenTelemetry:

#### HTTP Tracing
- All HTTP requests are automatically traced using the `otelmux` middleware
- Each incoming HTTP request creates a span
- Trace propagation headers are automatically handled

#### gRPC Tracing
- All gRPC calls are instrumented with interceptors
- Both unary and streaming calls are traced
- Trace context is propagated across service boundaries

## How to Use

### 1. Start the Services

```bash
docker-compose up -d
```

### 2. Access Jaeger UI

Open your browser and go to:
```
http://localhost:16686
```

### 3. Generate Traces

Make requests to your services through the API Gateway:

```bash
# Get all blogs (REST)
curl http://localhost:8080/api/v1/blogs

# Get blogs via gRPC (internal)
# This will be called by the API Gateway
```

### 4. View Traces

1. Go to Jaeger UI at `http://localhost:16686`
2. Select "blog-service" from the service dropdown
3. Click "Find Traces"
4. You will see all traces for the selected service

## What You'll See

### Trace Visualization

- **Service Map**: See how services interact with each other
- **Timeline View**: See the duration of each operation
- **Span Details**: View detailed information about each span
  - HTTP method and path
  - Request/response sizes
  - Database query information
  - gRPC call details

### Trace Context Propagation

Traces automatically propagate through:
- HTTP headers (traceparent, tracestate)
- gRPC metadata
- Service-to-service calls

## Implementation Details

### Files Added/Modified

1. **docker-compose.yml**: Added Jaeger service configuration
2. **services/blog-service/internal/tracing/tracing.go**: OpenTelemetry initialization
3. **services/blog-service/cmd/api/main.go**: Tracing setup and middleware
4. **services/blog-service/internal/grpc/grpc_server.go**: gRPC instrumentation
5. **services/blog-service/go.mod**: OpenTelemetry dependencies

### Dependencies

```go
go.opentelemetry.io/otel v1.27.0
go.opentelemetry.io/otel/exporters/jaeger v1.17.0
go.opentelemetry.io/otel/sdk v1.27.0
go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.52.0
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.52.0
```

## Troubleshooting

### No Traces Appearing?

1. Check if Jaeger is running:
   ```bash
   docker ps | grep jaeger
   ```

2. Check blog-service logs:
   ```bash
   docker logs soa-tourist-app-blog-service
   ```

3. Verify tracing initialization:
   Look for: "Starting blog service with tracing enabled"

### Jaeger UI Not Loading?

1. Check if port 16686 is accessible:
   ```bash
   curl http://localhost:16686
   ```

2. Check Jaeger container logs:
   ```bash
   docker logs soa-tourist-app-jaeger
   ```

## Future Enhancements

To add tracing to other services:

1. Add OpenTelemetry dependencies to `go.mod`
2. Create `internal/tracing/tracing.go` (same structure as blog-service)
3. Initialize tracing in `main.go`
4. Add middleware for HTTP or interceptors for gRPC
5. Update docker-compose.yml to add Jaeger dependency

## References

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Go SDK](https://pkg.go.dev/go.opentelemetry.io/otel)

