package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/francoispqt/onelog"
	"github.com/getsentry/sentry-go"
	sentryfasthttp "github.com/getsentry/sentry-go/fasthttp"
	"github.com/valyala/fasthttp"
)

var (
	Logger     *onelog.Logger
	config     Config
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.json", "Config file.")

	Logger = onelog.New(
		os.Stdout,
		onelog.ALL, // shortcut for onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
	)
}

func main() {
	flag.Parse()

	err := config.ParseConfigFile(configFile)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	handler := &Interceptor{
		clientHTTP: http.Client{},
	}

	fastHTTPhandler := handler.HandleFastHTTP
	if config.Sentry.Dsn != "" {
		err = sentry.Init(
			sentry.ClientOptions{
				Dsn:              config.Sentry.Dsn,
				Debug:            config.Sentry.Debug,
				AttachStacktrace: config.Sentry.AttachStacktrace,
				IgnoreErrors:     config.Sentry.IgnoreErrors,
				ServerName:       config.Sentry.ServerName,
				Release:          config.Sentry.Release,
				Environment:      config.Sentry.Environment,
				MaxBreadcrumbs:   config.Sentry.MaxBreadcrumbs,
			})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}

		sentryHandler := sentryfasthttp.New(sentryfasthttp.Options{})

		fastHTTPhandler = sentryHandler.Handle(fastHTTPhandler)
	}

	Logger.Info(fmt.Sprintf("Starting proxy server on %s", addr))
	if err := fasthttp.ListenAndServe(addr, fastHTTPhandler); err != nil {
		Logger.Fatal(fmt.Sprintf("ListenAndServe: %v", err))
	}
}
