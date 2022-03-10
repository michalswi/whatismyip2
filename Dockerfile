ARG GOLANG_VERSION
ARG ALPINE_VERSION

# build
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG APPNAME

RUN apk --no-cache add make git libc-dev libpcap-dev gcc

WORKDIR /app

COPY server server
COPY main.go main.go
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN make build

# execute
FROM alpine:${ALPINE_VERSION}

RUN apk --no-cache add libc-dev libpcap-dev gcc

ARG APPNAME
ENV SERVER_PORT 8080

COPY --from=builder /app/${APPNAME} /usr/bin/${APPNAME}

CMD ["whatismyip"]