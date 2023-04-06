// Package logger provides a logging utility for Go applications using the Echo web framework.
package logger

import (
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type ServiceName string
type GraylogEndpoint string

// GetLogID generate and return UUID.
func GetLogID() string {
	// Generate UUID
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "failed to generate UUID"
	}

	return uuid.String()
}

// Logger interface defines the Log method for logging messages.
type Logger interface {
	Log(c echo.Context) *logrus.Entry
	DefaultLog() *logrus.Entry
}

// LogImplementation struct implements the Logger interface.
type LogImplementation struct {
	Echo            *echo.Echo
	Logrus          *logrus.Logger
	ServiceName     string
	GraylogEndpoint string
}

// Log is a method of LogImplementation struct that logs the message at the given level and set some field from echo, like uri, remote ip & etc.
func (logger *LogImplementation) Log(c echo.Context) *logrus.Entry {
	// Create a log entry with fields for the HTTP request information.
	var email string
	session := c.Get("email")
	if session != nil {
		email = session.(string)
	}
	return logger.Logrus.WithFields(logrus.Fields{
		"uri":        c.Request().RequestURI,
		"remote_ip":  c.RealIP(),
		"host":       c.Request().Host,
		"method":     c.Request().Method,
		"error":      c.Error,
		"user_agent": c.Request().UserAgent(),
		"uri_path":   c.Path(),
		"user":       email,
		"log_id":     GetLogID(),
	})
}

// DefaultLog is method of LogImplementation struct with no additional field, DefaultLog just have graylog hook.
func (logger *LogImplementation) DefaultLog() *logrus.Entry {
	return logger.Logrus.WithFields(logrus.Fields{
		"log_id": GetLogID(),
	})
}

// NewLogger creates a new instance of Logger with Echo, Logrus, service name, and Graylog endpoint.
func NewLogger(e *echo.Echo, log *logrus.Logger, serviceName ServiceName, graylogEndpoint GraylogEndpoint) Logger {
	// Set log formatter to JSON
	log.SetFormatter(&logrus.JSONFormatter{})

	// Report the filename and line number of the calling function.
	log.SetReportCaller(true)

	// Add a hook to send logs to Graylog.
	hook := graylog.NewGraylogHook(string(graylogEndpoint), map[string]interface{}{"service": serviceName})
	log.AddHook(hook)

	// Middleware for logging all request inside Echo
	// Example output:
	// {"URI":"/test?oke=oke\u0026test=test","error":"code=404, message=Not Found","host":"localhost:8088","level":"info","method":"GET","msg":"request echo","query_param":null,"remote_ip":"::1","status":404,"time":"2023-03-08T21:03:43+07:00","uri_path":"/test","user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogError:     true,
		LogUserAgent: true,
		LogURIPath:   true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			// Log the HTTP request information in JSON format.
			log.WithFields(logrus.Fields{
				"uri":        values.URI,
				"status":     values.Status,
				"remote_ip":  values.RemoteIP,
				"host":       values.Host,
				"method":     values.Method,
				"error":      values.Error,
				"user_agent": values.UserAgent,
				"uri_path":   values.URIPath,
			}).Info("echo request")

			return nil
		},
	}))

	// Return a new LogImplementation object.
	return &LogImplementation{Echo: e, Logrus: log}
}
