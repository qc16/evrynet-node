FROM node:7

WORKDIR "/opt/explorer"
COPY ./web/ ./
COPY ./env_setup.sh .
RUN chmod +x env_setup.sh

RUN ./env_setup.sh
RUN npm install
CMD npm start

EXPOSE 8080
