FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Node ------"

ARG node_id
ARG has_metric
ENV ID $node_id
ENV HAS_METRIC $has_metric

WORKDIR "/node_$ID"
COPY ./gev ./
COPY ./genesis.json ./
COPY ./start_node.sh ./

RUN mkdir -pv ./log
RUN chmod +x ./start_node.sh

CMD ./start_node.sh
