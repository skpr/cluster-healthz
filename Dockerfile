FROM golang:1.13
RUN go get github.com/mitchellh/gox
ADD . /go/src/github.com/skpr/cluster-healthz
WORKDIR /go/src/github.com/skpr/cluster-healthz
RUN make build

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/skpr/cluster-healthz/bin/cluster-healthz_linux_amd64 /usr/local/bin/cluster-healthz
CMD ["cluster-healthz"]
