FROM golang:1.12-alpine as builder
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR "/evrynet"
COPY ./project ./project
RUN cd ./project && \
    go build ./cmd/bootnode

FROM alpine:latest
RUN apk add --no-cache ca-certificates

RUN echo "------ Bootnode ------"

WORKDIR "/bootnode"
COPY --from=builder ./evrynet/project/bootnode .

CMD ./bootnode -nodekeyhex "$NODEKEYHEX" -verbosity "$VERBOSITY" -addr ":$PORT"

EXPOSE $PORT