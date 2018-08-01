FROM golang:1.10-alpine AS build
WORKDIR /go/src/github.com/fankserver/spaceengineers-metrics
COPY . .
RUN apk add --no-cache alpine-sdk \
    && go get ./... \
    && go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN adduser -D -u 679 semetric
USER semetric

# Add app
COPY --from=build /go/src/github.com/fankserver/spaceengineers-metrics/app /app

# This container will be executable
ENTRYPOINT ["/app"]