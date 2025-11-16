package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/stannisl/pull-request-service/internal/config"
)

type Builder struct {
	port            string
	host            string
	handler         http.Handler
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

func NewBuilder() *Builder {
	return &Builder{
		host:            config.DefaultHTTPHost,
		port:            config.DefaultHTTPPort,
		readTimeout:     config.DefaultReadTimeout,
		writeTimeout:    config.DefaultWriteTimeout,
		shutdownTimeout: config.DefaultShutdownTimeout,
	}
}

func (b *Builder) WithHost(host string) *Builder {
	return b
}

func (b *Builder) WithPort(port string) *Builder {
	if port != "" {
		b.port = port
	}
	return b
}

func (b *Builder) WithHandler(handler http.Handler) *Builder {
	b.handler = handler
	return b
}

func (b *Builder) WithReadTimeout(timeout time.Duration) *Builder {
	if timeout > 0 {
		b.readTimeout = timeout
	}
	return b
}

func (b *Builder) WithWriteTimeout(timeout time.Duration) *Builder {
	if timeout > 0 {
		b.writeTimeout = timeout
	}
	return b
}

func (b *Builder) WithShutdownTimeout(timeout time.Duration) *Builder {
	if timeout > 0 {
		b.shutdownTimeout = timeout
	}
	return b
}

func (b *Builder) getAddress() string {
	return fmt.Sprintf("%s:%s", b.host, b.port)
}

func (b *Builder) Build() (*http.Server, time.Duration, error) {
	if b.handler == nil {
		return nil, b.shutdownTimeout, fmt.Errorf("handler is required to build a server")
	}

	return &http.Server{
		Addr:         b.getAddress(),
		Handler:      b.handler,
		ReadTimeout:  b.readTimeout,
		WriteTimeout: b.writeTimeout,
	}, b.shutdownTimeout, nil
}
