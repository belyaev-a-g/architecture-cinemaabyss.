FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# PROXY
#ENV HTTP_PROXY=http://172.16.0.250:4343
#ENV HTTPS_PROXY=http://172.16.0.250:4343

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o movies-service .

# Use a smaller image for the final application
FROM alpine:latest  

# PROXY
#ENV HTTP_PROXY=http://172.16.0.250:4343
#ENV HTTPS_PROXY=http://172.16.0.250:4343

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/movies-service .

# Expose the application port
EXPOSE 8081

# PROXY
#ENV HTTP_PROXY=
#ENV HTTPS_PROXY=

# Command to run the executable
CMD ["./movies-service"]
