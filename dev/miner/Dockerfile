# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-ethereum
RUN cd /go-ethereum && make geth

# -------------------------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR "/root"
RUN echo "------ Miner ------"

ARG password
ARG privatekey
RUN echo $password > ~/.accountpassword
RUN echo $privatekey > ~/.privatekey
ADD ./genesis.json ./genesis.json

COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

RUN geth init genesis.json
RUN geth account import --password ~/.accountpassword  ~/.privatekey

CMD exec geth --bootnodes "enode://$bootnodeId@$bootnodeIp:30301" --networkid 15 --verbosity 4  --syncmode "full" --mine --etherbase $address --minerthreads 2 --unlock $address --password ~/.accountpassword

EXPOSE 8545