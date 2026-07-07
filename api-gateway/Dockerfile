FROM golang:1.25-alpine AS build

WORKDIR /workspace

COPY contracts/go.mod contracts/go.sum ./contracts/

COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
WORKDIR /workspace/api-gateway
RUN go mod download

WORKDIR /workspace
COPY contracts/ ./contracts/
COPY api-gateway/ ./api-gateway/

WORKDIR /workspace/api-gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/srv ./cmd/api

FROM alpine:3.21 AS production

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone

COPY --from=build /bin/srv /app/srv

WORKDIR /app

EXPOSE 8080

USER nobody

ENTRYPOINT ["/app/srv"]
