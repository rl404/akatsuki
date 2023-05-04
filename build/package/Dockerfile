# Golang base image
FROM golang:1.20 as go_builder
LABEL stage="akatsuki_builder"
WORKDIR /akatsuki
COPY . .
WORKDIR /akatsuki/cmd/akatsuki
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod vendor -o akatsuki

# New stage from scratch
FROM alpine:3.16
LABEL stage="akatsuki_binary"
RUN apk add --no-cache ca-certificates && apk add --no-cache tzdata
COPY --from=go_builder /akatsuki/cmd/akatsuki/akatsuki /cmd/akatsuki/akatsuki
WORKDIR /cmd/akatsuki
