FROM golang:1.12-alpine as builder
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR "/evrynet"
COPY ./project ./project
RUN cd ./project && \
    go build ./cmd/gev


FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Node ------"

WORKDIR "/node"
COPY --from=builder ./evrynet/project/gev ./
COPY ./genesis.json ./
COPY ./start_node.sh ./
COPY ./stop_node.sh ./

RUN mkdir -pv ./log
RUN chmod +x ./start_node.sh

CMD bash ./start_node.sh