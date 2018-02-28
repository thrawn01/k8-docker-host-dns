# build the fix-host-dns binary
FROM golang:1.9 AS build-fix-host-dns

ENV GOPATH /go
RUN mkdir -p /go/src && mkdir -p /go/bin

WORKDIR /go/src/github.com/thrawn01/k8-docker-host-dns
COPY . .
RUN go install -ldflags "-linkmode external -extldflags -static" github.com/thrawn01/k8-docker-host-dns/cmd/fix-host-dns

# Final stage
FROM justincormack/nsenter1
WORKDIR /
COPY --from=build-fix-host-dns /go/bin/fix-host-dns /
ENTRYPOINT ["/fix-host-dns"]
