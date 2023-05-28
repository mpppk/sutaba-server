FROM golang:1-buster AS builder
ENV GO111MODULE on
RUN mkdir /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /src
WORKDIR /src
RUN make build

FROM alpine:3.18
RUN mkdir /lib64
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /src/sutaba-server /usr/local/bin
COPY --from=builder /src/config.yml /etc
ENTRYPOINT ["/usr/local/bin/sutaba-server", "start", "--config", "/etc/config.yml"]
