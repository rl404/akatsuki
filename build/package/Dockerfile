# Golang base image
FROM golang:1.27.4 AS go_builder
LABEL stage="akatsuki_builder"
WORKDIR /akatsuki
COPY . .
WORKDIR /akatsuki/cmd/akatsuki
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o akatsuki

# New stage from scratch
FROM alpine:3.24.4
LABEL stage="akatsuki_binary"
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 akatsuki
COPY --from=go_builder /akatsuki/cmd/akatsuki/akatsuki /usr/local/bin/akatsuki
WORKDIR /home/akatsuki
USER akatsuki
ENTRYPOINT ["/usr/local/bin/akatsuki"]
