FROM golang:1.12-alpine as builder
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY . /evrynet-client
WORKDIR "/evrynet-client"
RUN go build ./cmd/gev
RUN pwd
RUN ls

FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Node ------"

COPY --from=builder /evrynet-client/gev /
COPY ./deploy/devnet/node/genesis.json /
COPY ./dockerfiles/node/entrypoint.sh /

RUN chmod +x entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["/gev", "--help"]
