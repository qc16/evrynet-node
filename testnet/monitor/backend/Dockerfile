##FROM alpine:latest
#FROM ubuntu:xenial
#
#RUN apt-get update \
#      && apt-get install -y wget nodejs npm ntp git\
#      && rm -rf /var/lib/apt/lists/*
#
#RUN ln -s /usr/bin/nodejs /usr/bin/node
#RUN echo "------ Monitor Backend ------"
#
#WORKDIR "/opt"
#
#RUN git clone https://github.com/cubedro/eth-net-intelligence-api.git
#
#WORKDIR "/opt/eth-net-intelligence-api"
#
#RUN npm install -g pm2
#ADD ./app.json ./app.json
#RUN npm install
#
#CMD exec pm2-docker start app.json



#FROM ethereum/client-go
FROM alpine:latest
RUN echo "------ Monitor Backend ------"

RUN apk add --update git bash nodejs nodejs-npm perl

WORKDIR "/opt"
RUN git clone https://github.com/cubedro/eth-net-intelligence-api &&\
    cd eth-net-intelligence-api &&\
    npm install &&\
    npm install -g pm2

WORKDIR "/opt/eth-net-intelligence-api"

#RUN npm install -g pm2
ADD ./app.json ./app.json
#RUN npm install

CMD exec pm2-docker start ./app.json