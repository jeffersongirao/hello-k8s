FROM golang:1.10-alpine

WORKDIR /go/src/github.com/jeffersongirao/hello-k8s
COPY . .
RUN ./scripts/build.sh

FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates
COPY --from=0 /go/src/github.com/jeffersongirao/hello-k8s/bin/linux/hello-k8s /usr/local/bin

ENTRYPOINT ["/usr/local/bin/hello-k8s"]