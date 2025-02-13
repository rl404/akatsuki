# Golang base image
FROM golang:1.24 AS go_builder
LABEL stage="akatsuki_builder"
WORKDIR /akatsuki
COPY . .
WORKDIR /akatsuki/cmd/akatsuki
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o akatsuki

# New stage from scratch
FROM alpine:3.21
LABEL stage="akatsuki_binary"
RUN apk add --no-cache ca-certificates && apk add --no-cache tzdata
COPY --from=go_builder /akatsuki/cmd/akatsuki/akatsuki /cmd/akatsuki/akatsuki
WORKDIR /cmd/akatsuki
