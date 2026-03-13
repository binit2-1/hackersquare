# build go binary
FROM golang:1.22-alpine AS builder

# install git and ca-certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
# CGO_ENABLED=0 creates a statically linked binary (no C dependencies)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o hackersquare-api ./cmd/api


# Step 2: Create a tiny production image
FROM scratch

# Import the root ca-certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/hackersquare-api .

EXPOSE 8080


# Command to run the executable
CMD ["./hackersquare-api"]
