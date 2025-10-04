# syntax=docker/dockerfile:1

FROM golang:1.22.12-alpine AS builder

ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/server ./cmd/server
RUN go build -o /out/migrate ./cmd/migrate

FROM alpine:3.19

RUN adduser -S -D -H -h /app appuser \
 && apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /out/server /usr/local/bin/server
COPY --from=builder /out/migrate /usr/local/bin/migrate
COPY migrations ./migrations

USER appuser

ENV PORT=8080
EXPOSE 8080

CMD ["server"]
