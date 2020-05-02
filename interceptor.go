package main

import (
	"encoding/json"
	"fmt"
	"github.com/francoispqt/onelog"
	"github.com/francoispqt/onelog/log"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Interceptor struct{}

func (i *Interceptor) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	log.Info(fmt.Sprintf(`%s %s`, string(ctx.Method()), ctx.URI().String()))

	endpoint, ok := config.Router[string(ctx.Path())]
	if !ok {
		log.Error("Unsupported path")
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	u, err := url.Parse(endpoint.Url)
	if err != nil {
		log.Error(err.Error())
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}

	var body map[string]interface{}

	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		log.Error(err.Error())
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	// Filter
	if endpoint.Filter.Tags != nil {
		body["tags"] = filterTags(endpoint.Filter.Tags, body["tags"].([]interface{}))
	}

	if endpoint.Filter.Breadcrumbs != nil {
		body["breadcrumbs"] = filterBreadcrumbs(endpoint.Filter.Breadcrumbs, body["breadcrumbs"].(map[string]interface{}))
	}

	if endpoint.Filter.Extra != nil {
		body["extra"] = filterMap(endpoint.Filter.Extra, body["extra"].(map[string]interface{}))
	}

	bodyFiltered, _ := json.Marshal(body)

	// Forward the request filtered
	req := &ctx.Request
	res := &ctx.Response

	req.Header.SetHost(u.Host)
	req.SetRequestURI(endpoint.Url)
	req.SetBody(bodyFiltered)

	client := &fasthttp.HostClient{
		Addr: u.Host,
	}

	err = client.Do(req, res)
	if err != nil {
		ctx.Error("Server Error", fasthttp.StatusInternalServerError)
		log.Fatal(fmt.Sprintf("ServeHTTP: %v", err))
	}

	log.InfoWithFields("Response received", func(e onelog.Entry) {
		e.Int("status", res.StatusCode())
		e.String("host", res.RemoteAddr().String())
	})
}
