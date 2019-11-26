# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder
RUN echo "------ Node ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /evrynet
RUN cd /evrynet && go build ./cmd/gev

# -------------------------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR "/root"
COPY --from=builder /evrynet/gev /usr/local/bin/
COPY --from=builder /evrynet/dev/onenode/nodedata ./nodedata

ENV password "123"
ENV privatekey "ce900e4057ef7253ce737dccf3979ec4e74a19d595e8cc30c6c5ea92dfdd37f1"
ENV address "0x560089ab68dc224b250f9588b3db540d87a66b7a"

RUN echo $password > /root/.accountpassword
RUN echo $privatekey > /root/.privatekey
RUN chmod +x nodedata/one_node_genesis.json

CMD gev --datadir nodedata/data init nodedata/one_node_genesis.json && \
    gev account import --datadir nodedata/data --password /root/.accountpassword /root/.privatekey && \
    gev --datadir nodedata/data --nodiscover --tendermint.blockperiod 1 --gasprice 1000000000 --syncmode full --mine --networkid 15 \
    --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --port 30303 \
    --password /root/.accountpassword --unlock $address --etherbase $address \
    --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 --allow-insecure-unlock

EXPOSE 8545
EXPOSE 30303
