FROM golang:1.13.4-alpine AS build

WORKDIR /go/src/riffraff
COPY . /go/src/riffraff

RUN apk add --no-cache make

RUN make clean build

FROM alpine

ARG VERSION=dev
ARG COMMIT=00000000000000000000

LABEL maintainer="Adam Veldhousen <adam@vdhsn.com>"
LABEL version=${VERSION}
LABEL commit=${COMMIT}

WORKDIR /usr/local/bin
COPY --from=build /go/src/riffraff/.bin/riffraff /usr/local/bin/

RUN apk add --no-cache ca-certificates
RUN adduser -u 1000 -D -s /bin/ash riffraff \
    && chown -R riffraff:riffraff /usr/local/bin/riffraff

USER riffraff

EXPOSE 8080

VOLUME ["/data"]

ENTRYPOINT [ "/usr/local/bin/riffraff", "-bind", "0.0.0.0",  "-port", "8080", "-accesslog", "-data", "/data/data.json"]