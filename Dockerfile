FROM golang:1.24.2-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth ./cmd/app/app.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

COPY --from=builder --chown=appuser:appgroup /app/auth /usr/local/bin/auth

RUN mkdir -p /app/logs && chown -R appuser:appgroup /app/logs

RUN mkdir -p /configs && chown -R appuser:appgroup /configs

USER appuser

EXPOSE 50551

CMD ["/usr/local/bin/auth"]