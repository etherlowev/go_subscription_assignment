FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY internal/ ./internal
COPY cmd/ ./cmd
COPY docs/ ./docs
RUN /go/bin/swag init -g cmd/main.go --parseDependency

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd


FROM alpine:latest

RUN adduser -D appuser && mkdir /app

WORKDIR /app

COPY --chown=appuser:appuser --from=builder /app/main .
COPY --chown=appuser:appuser config.yaml .

USER appuser

EXPOSE 8080

CMD ["./main"]