## Building binaries ##
FROM golang:1.17-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /tmp/go-stream-broadcast

COPY . .

RUN go mod download
RUN make build_reader build_broadcaster

## Creating final image ##
FROM alpine:latest 

RUN apk add ca-certificates

COPY --from=builder /tmp/go-stream-broadcast/dist/reader /app/reader
COPY --from=builder /tmp/go-stream-broadcast/dist/broadcaster /app/broadcaster

CMD ["/app/reader"]