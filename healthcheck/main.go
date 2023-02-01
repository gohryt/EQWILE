package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gohryt/EQWILE/healthcheck/checker"

	"github.com/pelletier/go-toml/v2"
	"github.com/valyala/fasthttp"
)

type (
	Configuration struct {
		Name string `toml:"name"`
		Host string `toml:"host"`
		Port int    `toml:"port"`

		Checker checker.Configuration
	}
)

const (
	HSTSMaxAge int = 31536000
)

func main() {
	path := ".configuration"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	configuration := &Configuration{}

	err = toml.NewDecoder(file).Decode(configuration)
	if err != nil {
		log.Fatal(err)
	}

	err = Main(configuration)
	if err != nil {
		log.Fatal(err)
	}
}

func Main(configuration *Configuration) (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ctx, cancelCause := context.WithCancelCause(ctx)
	defer cancelCause(nil)

	checker := checker.Constructor(&configuration.Checker)

	checker.Register("status_code", func(response *fasthttp.Response) any {
		if response.StatusCode() == fasthttp.StatusOK {
			return true
		}

		return false
	})

	server := fasthttp.Server{
		Name:    configuration.Host,
		Handler: func(ctx *fasthttp.RequestCtx) { ctx.WriteString("Hello EQWILE") },

		Concurrency:  1024 * 16,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  16 * time.Second,
		TCPKeepalive: true,

		KeepHijackedConns: true,
		StreamRequestBody: true,
		CloseOnShutdown:   true,
	}
	defer server.Shutdown()

	http, err := net.Listen("tcp", (":" + strconv.Itoa(configuration.Port)))
	if err != nil {
		return
	}

	go func() {
		cancelCause(checker.Run(ctx))
	}()

	go func() {
		cancelCause(server.Serve(http))
	}()

	select {
	case <-ctx.Done():
		err = context.Cause(ctx)
		if err == context.Canceled {
			err = nil
		}
	}

	return
}
