FROM alpine:latest
RUN apk add --no-cache ca-certificates

RUN echo "------ Bootnode ------"

WORKDIR "/bootnode"
COPY ./bin/bootnode .

CMD ./bootnode -nodekeyhex "$NODE_HEX_KEY" -verbosity 9

EXPOSE 30301/udp
EXPOSE 30303/udp