FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /watcher .

FROM alpine:3.23

RUN apk add --no-cache ca-certificates

COPY --from=builder /watcher /watcher

USER 1000:1000

ENTRYPOINT ["/watcher"]
