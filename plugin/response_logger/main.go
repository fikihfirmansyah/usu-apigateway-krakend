package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/proxy"
	"go.uber.org/zap"
)

// ResponseLoggerMiddleware is the main structure for the plugin
type ResponseLoggerMiddleware struct {
	logger *zap.Logger
}

// New creates a new instance of ResponseLoggerMiddleware
func New(logger *zap.Logger) *ResponseLoggerMiddleware {
	return &ResponseLoggerMiddleware{
		logger: logger,
	}
}

// NewBackendFactory is the factory function for the plugin
func (m *ResponseLoggerMiddleware) NewBackendFactory(next proxy.BackendFactory) proxy.BackendFactory {
	return func(cfg *config.Backend) proxy.Proxy {
		next := next(cfg)
		return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
			resp, err := next(ctx, req)
			if err != nil {
				return resp, err
			}

			// Log the response
			m.logResponse(req, resp)

			return resp, nil
		}
	}
}

// logResponse logs the response details
func (m *ResponseLoggerMiddleware) logResponse(req *proxy.Request, resp *proxy.Response) {
	// Extract relevant information from the response
	statusCode := resp.Metadata.StatusCode
	responseBody, _ := ioutil.ReadAll(resp.Io)

	// Create a log entry
	logEntry := map[string]interface{}{
		"timestamp":     time.Now().Format(time.RFC3339),
		"method":        req.Method,
		"path":          req.Path,
		"status_code":   statusCode,
		"response_body": string(responseBody),
	}

	// Convert log entry to JSON
	jsonLog, err := json.Marshal(logEntry)
	if err != nil {
		m.logger.Error("Failed to marshal log entry", zap.Error(err))
		return
	}

	// Log the entry
	m.logger.Info(string(jsonLog))
}

// BackendFactory is the exported symbol for registering the plugin
var BackendFactory = func(next proxy.BackendFactory) proxy.BackendFactory {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return New(logger).NewBackendFactory(next)
}

func main() {}