FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Miners ------"

WORKDIR "/node"
COPY ./testnet/miner/bin/gev .
COPY ./testnet/miner/bin/makegenesis .
COPY ./testnet/miner/utils/ .

ARG numberOfMiners
ENV NumberOfMiners $numberOfMiners
RUN chmod +x *.sh

CMD GEV=./gev BOOTNODE=./bootnode bash ./gev_cluster.sh ./testnet 15 $NumberOfMiners \
    --bootnodes "enode://$BOOTNODE_ID@$BOOTNODE_IP:30301" && \
    tail -f /dev/null # Keep container run

