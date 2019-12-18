# run 2 dockers
- access to the dev directory and run the commands here
```
$ docker volume create influxdb-volume
$ docker volume create grafana-volume
$ docker-compose up -d geth-grafana geth-influxdb
```
# config node to sync the data metrics to influx db
- run node with parameter 

```
--metrics --metrics.influxdb
```

# import dashboard
- navigate to grafana via url http://localhost:3001
- login to admin site
- config to connect to influxdb with parameters here:
  - endpoint: localhost:8086
  - database: geth
  - user: test
  - password: test
- create dashboard via import method
