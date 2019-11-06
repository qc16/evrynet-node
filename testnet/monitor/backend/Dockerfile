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