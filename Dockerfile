FROM golang:1.18.3-alpine AS builder
LABEL description="GT - Geo Tracking" version="1.0"

RUN apk --no-cache add --update bash
RUN apk --no-cache add --update alpine-sdk

RUN mkdir -p /app
ADD . /app
WORKDIR /app

RUN make install

RUN rm -rf /app

FROM golang:1.18.3-alpine
COPY --from=builder /go/bin/geo-tracking /go/bin/geo-tracking