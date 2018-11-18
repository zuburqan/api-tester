API Tester is a tool written in golang. It's purpose is to performance test any API.

It exposes prometheus metrics are to be used to analyse APIs performance

### Metrics Exposed:

1. Total requests split by method, code and API endpoint
2. Failed requests split by method code and endpoint
3. Total journeys split by journey name
4. Failed journeys split by name
5. Request duration split by different time buckets (0.25, 0.5, 1, 2.5, 5 and 10 seconds)
6. Journey duration split by different time buckets (0.25, 0.5, 1, 2.5, 5 and 10 seconds)

### Configuration:

api-tester.conf : This is used to specify the following properties: 

1. destination: This will be used in metric names as a prefix & is for identification purposes.
2. auth: auth type for the API. Currently "no_auth", "basic" and "digest" auths are supported
3. client_timeout: default 60 seconds. This is the golang HTTP client timeout.
4. sleep: number of seconds the API tester waits before repeating the journey set for a user.
5. users: the number of users which will execute the journeys concurrently.
6. log_level: "debug", "info" and "warn" levels available.
7. stats_host: the prometheus metrics server host IP, that will expose these metrics to be collected
8. stats_port: the prometheus metrics server host port, that will expose these metrics to be collected
9. journeys: the journeys to be executed, see example journey set for schema format.

api-tester-connection.conf: This is used to specify the following properties:

1. host: the API host along with port. example: http://foobar.com:8080
2. username: username of API
3. password: password of API

The following is an example journey set for running against Transport for London (TFL) API

https://github.com/zuburqan/api-tester/blob/master/etc/api-tester.conf.example

### Dependencies

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
