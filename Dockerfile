# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git curl

WORKDIR /evrynet-node
ADD . .

# Load all project dependencies
RUN ./vendor.sh

# Install golangci-lint tool
RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
				sh -s -- -b $GOPATH/bin v1.21.0

# Build the gev binary
#TODO: not sure why but if replace this by 'go run build/ci.go install', Jenkins CI is failed by out of memory
RUN go run build/ci.go install ./cmd/gev

RUN go build ./cmd/gev
RUN go build ./cmd/bootnode

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /evrynet-node/gev /usr/local/bin/
COPY --from=builder /evrynet-node/bootnode /usr/local/bin/

#--rpcport 8545
#--wsport 8546
#--graphql.port 8547
#--port 30303
EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["gev"]
CMD ["--rpc", "--rpcaddr","0.0.0.0", "--rpcvhosts", "*"]