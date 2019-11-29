FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Node ------"

WORKDIR "/node"
COPY ./bin/gev ./
COPY ./bin/genesis.json ./
COPY ./start_node.sh ./
COPY ./stop_node.sh ./

RUN mkdir -pv ./log
RUN chmod +x ./start_node.sh

CMD bash ./start_node.sh
