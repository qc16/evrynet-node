FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Miners ------"

WORKDIR "/node"
COPY ./testnet/miner/bin/gev .
COPY ./testnet/miner/bin/makegenesis .
COPY ./testnet/miner/utils/ .
#ADD ./testnet/genesis.json ./genesis.json

RUN chmod +x *.sh

CMD GETH=./gev bash ./gev_cluster.sh ./testnet 15 5 \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30301" && \
    tail -f /dev/null # Keep container run

