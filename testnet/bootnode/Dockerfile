FROM alpine:latest
RUN apk add --no-cache ca-certificates

RUN echo "------ Bootnode ------"

WORKDIR "/bootnode"
COPY ./testnet/bootnode/bin/bootnode .

CMD ./bootnode -nodekeyhex $NODE_HEX_KEY

EXPOSE 30301/udp
EXPOSE 30303/udp