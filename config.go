package main

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/orandin/go-ruler"
	"io/ioutil"
	"log"
)

type Config struct {
	Host   string       `json:"host"`
	Port   int64        `json:"port"`
	Sentry SentryConfig `json:"sentry"`
	Router Router       `json:"router"`
}

type Router map[string]Endpoint

type Endpoint struct {
	Url    string       `json:"url"`
	Filter SentryFilter `json:"filter"`
}

type SentryFilter struct {
	Tags        []SentryFilterRules `json:"tags"`
	Breadcrumbs []SentryFilterRules `json:"breadcrumbs"`
	Extra       []SentryFilterRules `json:"extra"`
}

type SentryFilterRules struct {
	Conditions []*ruler.Rule `json:"conditions"`
}

type SentryConfig struct {
	sentry.ClientOptions
	Dsn              string   `json:"dsn"`
	Debug            bool     `json:"debug"`
	AttachStacktrace bool     `json:"attach_stacktrace"`
	IgnoreErrors     []string `json:"ignore_errors"`
	ServerName       string   `json:"server_name"`
	Release          string   `json:"release"`
	Environment      string   `json:"environment"`
	MaxBreadcrumbs   int      `json:"max_breadcrumbs"`
}

func (c *Config) ParseConfigFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return nil
}
