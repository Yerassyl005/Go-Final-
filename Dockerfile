FROM golang:1.25-alpine AS deps

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

FROM deps AS builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/smartqueue ./cmd/api

FROM alpine:3.22 AS runtime

RUN addgroup -S smartqueue && adduser -S smartqueue -G smartqueue

WORKDIR /app

COPY --from=builder /out/smartqueue /app/smartqueue

USER smartqueue

EXPOSE 8080

ENTRYPOINT ["/app/smartqueue"]
