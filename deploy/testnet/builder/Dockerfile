# Build Geth in a stock Go builder container
FROM golang:1.12-alpine
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR "/evrynet"
COPY ./project ./project
RUN ls -la && \
    cd ./project && \
    go build ./cmd/gev && \
    go build ./cmd/bootnode && \
    ls -la && \
    cd .. && \
    cp ./project/gev ./ && \
    cp ./project/bootnode ./ && \
    rm -rf ./project
