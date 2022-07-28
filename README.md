# go-parse-log-uri
Simple Go utility to parse a Nginx log files, to extract all the layers in the WMS queries from log entries.

WMS is a standard interface for requesting geospatial map images.
There are several GeoSpatial tools  
[Mapserver](https://mapserver.org/ogc/wms_server.html)

## how to use

    goParseLog yourNginxCombinedAccess.log


## how to build

    go build -o goParseLog go-parse-log.go


## how to run

    go run go-parse-log.go

