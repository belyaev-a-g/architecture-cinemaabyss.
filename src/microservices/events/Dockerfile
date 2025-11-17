# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# PROXY
ENV HTTP_PROXY=http://172.16.0.250:4343
ENV HTTPS_PROXY=http://172.16.0.250:4343

# Initialize module and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# PROXY
ENV HTTP_PROXY=http://172.16.0.250:4343
ENV HTTPS_PROXY=http://172.16.0.250:4343

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8082

# Command to run
CMD ["./main"]
