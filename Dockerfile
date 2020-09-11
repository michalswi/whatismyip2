ARG GOLANG_VERSION
ARG ALPINE_VERSION

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION}

ARG APPNAME
ENV SERVER_ADDR "8080"

RUN apk --no-cache add make git libc-dev libpcap-dev gcc

WORKDIR /app

COPY whatismyip.go whatismyip.go
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN make build

CMD ["./whatismyip"]
# COPY ${APPNAME} /usr/bin/${APPNAME}
# CMD ["whatismyip"]