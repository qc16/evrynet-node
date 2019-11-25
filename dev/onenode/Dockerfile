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
COPY --from=builder /evrynet/tests ./
RUN chmod +x test_nodes/onenode/one_node_genesis.json

CMD gev --datadir test_nodes/onenode/data init test_nodes/onenode/one_node_genesis.json && \
    gev --datadir test_nodes/onenode/data --nodiscover --tendermint.blockperiod 1 --gasprice 1000000000 --syncmode full  --mine --networkid 15 \
    --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --port 30303 \
    --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3 --allow-insecure-unlock

EXPOSE 8545
EXPOSE 30303
