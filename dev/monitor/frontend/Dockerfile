FROM node:8-alpine

RUN apk add --update git
RUN echo "------ Monitor Frontend ------"

RUN git clone https://github.com/cubedro/eth-netstats.git

WORKDIR /eth-netstats

RUN npm install
RUN npm install -g grunt-cli
RUN grunt

EXPOSE 3000

CMD WS_SECRET=eth-net-stats-secret npm start