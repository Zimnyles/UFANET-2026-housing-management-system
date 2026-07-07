FROM golang:1.25-alpine AS build
WORKDIR /workspace
COPY contracts/go.mod contracts/go.sum ./contracts/
COPY notification-microservice/go.mod notification-microservice/go.sum ./notification-microservice/
WORKDIR /workspace/notification-microservice
RUN go mod download
WORKDIR /workspace
COPY contracts/ ./contracts/
COPY notification-microservice/ ./notification-microservice/
WORKDIR /workspace/notification-microservice
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/srv ./cmd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /bin/srv /app/srv
WORKDIR /app
EXPOSE 50055
USER nobody
ENTRYPOINT ["/app/srv"]
