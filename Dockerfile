# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /evrynet-node
RUN cd /evrynet-node && go build ./cmd/gev && go build ./cmd/bootnode && go build ./cmd/puppeth

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /evrynet-node/gev /usr/local/bin/

#--rpcport 8545
#--wsport 8546
#--graphql.port 8547
#--port 30303
EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT ["gev"]
CMD ["--rpc", "--rpcaddr","0.0.0.0", "--rpcvhosts", "*"]