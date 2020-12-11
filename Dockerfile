FROM golang:1.15-alpine as server_build

# Add build deps
RUN apk add --update gcc g++ git

COPY go.mod go.sum /appbuild/

COPY ./ /appbuild

RUN set -ex \
    && go version \
    && cd /appbuild \
    && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -o server

# Build deployable server
FROM alpine:latest

WORKDIR /opt/server

RUN set -ex \
    && apk add --update --no-cache ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

COPY --from=server_build /appbuild/server /opt/server

EXPOSE 8080

CMD ["./server"]
