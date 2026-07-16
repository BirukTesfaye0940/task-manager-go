# Step 1: Build the Go binary
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o task-manager-go main.go

# Step 2: Create a minimal, secure runner image
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Copy compiled binary from builder
COPY --from=builder /app/task-manager-go /task-manager-go

# Expose server port
EXPOSE 8080

# Run the app
ENTRYPOINT ["/task-manager-go"]
