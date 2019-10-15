# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR "/evrynet"
ADD . .
RUN go build ./cmd/gev && go build ./cmd/bootnode && go build -o ./makegenesis ./testnet/utils/makegenesis.go
