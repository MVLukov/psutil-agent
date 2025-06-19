# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files and download deps first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go app
RUN go build -o main .

# Final image
FROM ubuntu:latest

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/main .

# Copy templates and static assets
COPY --from=builder /app/templates/ templates/
COPY --from=builder /app/static/ static/

# Expose port (if needed)
EXPOSE 3000


# Run the app
CMD ["./main"]
