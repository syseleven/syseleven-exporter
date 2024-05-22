FROM golang:1.21-alpine3.19 as build

RUN apk add --no-cache --update git make

RUN mkdir /build
WORKDIR /build
COPY . .
RUN make build


FROM alpine:3.20

ARG REVISION
ARG VERSION

LABEL maintainer="SysEleven GmbH"
LABEL git.url="https://github.com/syseleven/syseleven-exporter"

RUN apk add --no-cache --update curl ca-certificates
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD curl --fail http://localhost:8080/_health || exit 1

USER nobody

COPY --from=build /build/bin/syselevenexporter /bin/syselevenexporter
EXPOSE 8080

ENTRYPOINT  [ "/bin/syselevenexporter" ]
