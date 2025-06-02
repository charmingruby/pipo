# Stage 1: Dependencies
FROM golang:1.23.2-alpine AS deps

WORKDIR /app

COPY service/api/go.mod service/api/go.sum ./
COPY lib/go.mod lib/go.sum ./lib/

RUN mkdir -p /go/src/github.com/charmingruby/pipo && \
    cp -r /app/lib /go/src/github.com/charmingruby/pipo/

# Workaround for the module path
# In reality, the module path should be a separate repository
# and should be imported as a dependency
RUN sed -i 's|replace github.com/charmingruby/pipo-lib => ../../lib|replace github.com/charmingruby/pipo-lib => /go/src/github.com/charmingruby/pipo-lib|' go.mod

RUN go mod download

# Stage 2: Build
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg

# Copy both api and lib source code
COPY service/api .
COPY lib ./lib

RUN mkdir -p /go/src/github.com/charmingruby/pipo && \
    cp -r /app/lib /go/src/github.com/charmingruby/pipo/

# Workaround for the module path
# In reality, the module path should be a separate repository
# and should be imported as a dependency
RUN sed -i 's|replace github.com/charmingruby/pipo-lib => ../../lib|replace github.com/charmingruby/pipo-lib => /go/src/github.com/charmingruby/pipo-lib|' go.mod

RUN go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Final stage: Run the application
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache shadow

# Create a non-root user and group
RUN groupadd -r appuser && useradd -r -g appuser appuser

COPY --from=builder /app/main .

# Change ownership of the application binary
RUN chown appuser:appuser /app/main

# Switch to the non-root user
USER appuser

CMD ["./main"]