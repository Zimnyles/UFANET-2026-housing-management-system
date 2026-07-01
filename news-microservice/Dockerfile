FROM golang:1.25-alpine AS build

WORKDIR /workspace

COPY contracts/go.mod contracts/go.sum ./contracts/

COPY news-microservice/go.mod news-microservice/go.sum ./news-microservice/
WORKDIR /workspace/news-microservice
RUN go mod download

WORKDIR /workspace
COPY contracts/ ./contracts/
COPY news-microservice/ ./news-microservice/

WORKDIR /workspace/news-microservice
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/srv ./cmd

FROM alpine:3.21 AS production

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone

COPY --from=build /bin/srv /app/srv

WORKDIR /app

EXPOSE 50053

USER nobody

ENTRYPOINT ["/app/srv"]
