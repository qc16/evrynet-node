# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder
RUN echo "------ Miner ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-ethereum
RUN cd /go-ethereum && make geth

# -------------------------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR "/root"

ARG password
ARG privatekey
ARG needRestore
ENV NEED_RESTORE $needRestore
RUN echo $password > /root/.accountpassword
RUN echo $privatekey > /root/.privatekey
ADD ./genesis.json /root/genesis.json

COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

COPY ./dev/miner/restore.sh restore.sh
RUN chmod +x restore.sh

CMD ./restore.sh && \
    geth --networkid 15 --verbosity 4 --syncmode "fast" \
    --bootnodes "enode://$bootnodeId@$bootnodeIp:30301" \
    --etherbase $address --mine --minerthreads 2

EXPOSE 8545