FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

# PROXY
ENV HTTP_PROXY=http://172.16.0.250:4343
ENV HTTPS_PROXY=http://172.16.0.250:4343

RUN go mod tidy && go build -o events-service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/events-service .

EXPOSE 8082

ENV PORT=8082

# PROXY
ENV HTTP_PROXY=
ENV HTTPS_PROXY=

CMD ["./events-service"] 
