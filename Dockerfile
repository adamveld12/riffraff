FROM golang:1.13.3-alpine AS build

WORKDIR /go/src/riffraff
COPY . /go/src/riffraff

RUN apk add --no-cache make

RUN make clean build

FROM alpine

ARG VERSION=dev

LABEL maintainer="Adam Veldhousen <adam@vdhsn.com>"
LABEL version=${VERSION}

WORKDIR /usr/local/bin
COPY --from=build /go/src/riffraff/.bin/riffraff /usr/local/bin/

RUN apk add --no-cache ca-certificates
RUN adduser -u 1000 -D -s /bin/ash riffraff \
    && chown -R riffraff:riffraff /usr/local/bin/riffraff

USER riffraff

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/riffraff", "-port", "8080", "-accesslog"]