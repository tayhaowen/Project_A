# Multi-stage Dockerfile for the Project A Go service

# ---- Build stage ----
ARG GO_TOOLCHAIN=go1.25.4
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Ensure the requested toolchain version is downloaded/used when running go commands
ENV GOTOOLCHAIN=${GO_TOOLCHAIN}

# Install build dependencies
RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary (and optionally .env if you really need it inside the image)
COPY --from=builder /app/server /app/server

# Expose the HTTP port
EXPOSE 8080

# Run as non-root user provided by distroless base image
USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
