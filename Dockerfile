# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git curl

WORKDIR /go-ethereum

ADD . /go-ethereum

# Load all project dependencies
RUN ./vendor.sh

# Install golangci-lint tool
RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
				sh -s -- -b $GOPATH/bin v1.21.0

# Build the gev binary
RUN go run build/ci.go install

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ethereum/build/bin/gev /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["gev"]
