# Build context: project root  (docker-compose: context: .)
FROM golang:1.25-alpine AS build

WORKDIR /workspace

# contracts (нужен для replace директивы)
COPY contracts/go.mod contracts/go.sum ./contracts/

# profile-microservice зависимости
COPY profile-microservice/go.mod profile-microservice/go.sum ./profile-microservice/
WORKDIR /workspace/profile-microservice
RUN go mod download

# исходники
WORKDIR /workspace
COPY contracts/ ./contracts/
COPY profile-microservice/ ./profile-microservice/

WORKDIR /workspace/profile-microservice
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/srv ./cmd

# ---

FROM alpine:3.21 AS production

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone

COPY --from=build /bin/srv /app/srv

WORKDIR /app

EXPOSE 50052

USER nobody

ENTRYPOINT ["/app/srv"]
