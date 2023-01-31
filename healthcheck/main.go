package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/valyala/fasthttp"
)

type (
	Configuration struct {
		Name string `toml:"name"`
		Host string `toml:"host"`

		URLs []URL
	}

	URL struct {
		URL        string
		CheckList  []string `toml:"checkList"`
		CheckCount int      `toml:"checkCount"`
	}
)

const (
	HSTSMaxAge int = 31536000
)

func main() {
	file, err := os.OpenFile(".configuration", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	configuration := &Configuration{}

	err = toml.NewDecoder(file).Decode(configuration)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(configuration)

	err = Main(configuration)
	if err != nil {
		log.Fatal(err)
	}
}

func Main(configuration *Configuration) (err error) {
	http, err := net.Listen("tcp", ":80")
	if err != nil {
		return
	}

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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	errs := make(chan error)

	go func() {
		errs <- server.Serve(http)
	}()

	select {
	case err = <-errs:
	case <-signals:
	}

	return
}
