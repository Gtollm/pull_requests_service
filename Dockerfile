FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pr-reviewer ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /root/

COPY --from=builder /app/pr-reviewer .

RUN chown appuser:appuser pr-reviewer

USER appuser

EXPOSE 8080

CMD ["./pr-reviewer"]