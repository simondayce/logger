# Logging Utility for Go Applications using Echo Web Framework

This package provides a logging utility for Go applications using the Echo web framework. The logger logs requests and responses, as well as other events, and sends them to a Graylog server for centralized logging.

## Getting Started
### Prerequisites
Before using the logger, you should have Go installed.

### Installation
To install the logger, run the following command:

```shell
go get github.com/SimondayCE/logger
```

### Environment variable
Make sure to set the environment variables serviceName and graylogEndpoint on your machine before running the application. The serviceName should be set to the name of your service, and the graylogEndpoint should be set to the IP address and port of your Graylog server.

Example:
- serviceName=auth
- graylogEndpoint=10.10.10.107:5000

### Usage
To use this package, import it in your Go code:

```go
import "github.com/SimondayCE/logger"
```

Then, create a new logger using `NewLogger` function:
```go
e := echo.New()
log := logrus.New()
logger := logger.NewLogger(e, log, os.Getenv("serviceName"), os.Getenv("graylogEndpoint"))
```

`NewLogger` function takes four parameters:

- `e`: an instance of `echo.Echo` struct.
- `log`: an instance of `logrus.Logger` struct.
- `serviceName`: a string containing the name of your service.
- `graylogEndpoint`: a string containing the endpoint of your Graylog server.

You can now use `logger.Log` function to log messages, like default logrus logger:
```go
e.GET("/", func(c echo.Context) error {
    log := logger.Log(c)
	log.Info("Log message here!")
    return c.String(http.StatusOK, "Hello, World!")
})
```

`logger.Log` function takes one parameter:

- `c`: an instance of `echo.Context` struct.

And `logger.Log` will return `*logrus.Entry`

You can use this logger module in the delivery layer, so you can pass `echo.Context` to `logger.Log` for get uri, remote ip & etc. Make sure to also pass `logger.Log` to the next layer, such as the use case and repository layers.

## Built With
- [Echo](https://github.com/labstack/echo) - A fast and unfancy micro web framework for Go.
- [Logrus](https://github.com/sirupsen/logrus) - A structured logger for Go.
- [logrus-graylog-hook](https://github.com/gemnasium/logrus-graylog-hook) - A hook for Logrus that sends logs to Graylog.