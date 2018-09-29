# spaceengineers-metrics

```
Usage of spaceengineers-metrics:
  -host string
    	host url of the rcon server (default "http://localhost:8080")
  -influxdb string
    	influxdb database (default "spaceengineers")
  -influxhost string
    	influxdb host (default "http://localhost:8086")
  -influxpass string
    	influx password
  -influxuser string
    	influx username
  -key string
    	rcon key
```

Docker example:
```
docker run -d --name spaceengineers-metrics fankserver/spaceengineers-metrics \
    -host http://localhost:8080 \
    -key foobar123 \
    -influxhost http://localhost:8086 \
    -influxuser asd
```
