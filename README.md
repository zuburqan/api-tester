API Tester is a tool written in golang. It's purpose is to performance test any API.

It exposes prometheus metrics which can be converted to graphite using the carboniser.

## Metrics Exposed:

```shell
Total requests split by method, code and API endpoint
Failed requests split by method code and endpoint
Total journeys split by journey name
Failed journeys split by name
Request duration split by different time buckets (0.25, 0.5, 1, 2.5, 5 and 10 seconds)
Journey duration split by different time buckets (0.25, 0.5, 1, 2.5, 5 and 10 seconds)
```

## Configuration:

api-tester.conf : This is used to specify the following properties: 
```
destination: set it to any name you like, usually the name of the API. This will be used in metric names as a prefix
auth: specify the auth type for the API. Currently "no_auth", "basic" and "digest" auths are supported
client_timeout: default 60 seconds. This is the golang HTTP client timeout.
sleep: how long a user of the API waits before starting the journey set again (in seconds)
users: the number of users which will execute the journeys concurrently
log_level: defaults to "warn" where you will only see errors if any. change to "info" to see all requests being executed and to "debug" to see full response body and status code
stats_host: the prometheus metrics server host IP, that will expose these metrics to be collected
stats_port: the prometheus metrics server host port, that will expose these metrics to be collected
journeys: the journeys to be executed. This has the name, any setup requests, the actual requests and any cleanup requests if required
```

api-tester-connection.conf: This is used to specify the following properties:

```
host: the API host along with port. example: http://foobar.com:8080
username: username of API
password: password of API
```

The following is an example journey set for running against Transport for London (TFL) API

https://github.com/zuburqan/api-tester/blob/master/etc/api-tester.conf


## Dependencies

*   [Go 1.10](https://golang.org/)
*   [Dep](https://github.com/golang/dep)

```shell
brew install go dep
```

## Installation

```shell
go get github.com/zuburqan/api-tester
cd $GOPATH/src/github.com/zuburqan/api-tester
dep ensure
```

## Running

```shell
go run main.go
```
