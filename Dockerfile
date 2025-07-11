# syntax=docker/dockerfile:1.5
FROM golang:1.25rc2-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o server ./cmd/hn

FROM node:20 AS frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Use a minimal image for the final stage
FROM debian:bookworm-slim AS final
WORKDIR /app
# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
# Copy the binary
COPY --from=builder /app/server .
COPY --from=frontend-builder /frontend/dist ./static
EXPOSE 30000
CMD ["./server"]