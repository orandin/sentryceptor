package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/francoispqt/onelog"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type Interceptor struct {
	clientHTTP http.Client
}

func (i *Interceptor) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	requestID := uuid.New().String()
	log := Logger.With(func(e onelog.Entry) {
		e.String("time", time.Now().Format(time.RFC3339))
		e.String("internalRequestID", requestID)
	})

	log.InfoWithFields("Sentry event received", func(e onelog.Entry) {
		e.String("method", string(ctx.Method()))
		e.String("path", string(ctx.URI().Path()))
	})

	endpoint, ok := config.Router[string(ctx.Path())]
	if !ok {
		log.Error("Unsupported path")
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	dsn, err := sentry.NewDsn(endpoint.Url)
	if err != nil {
		log.Error(err.Error())
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}

	url := dsn.StoreAPIURL()

	// Get body
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

	// Send filtered event to Sentry
	request, _ := http.NewRequest(
		http.MethodPost,
		url.String(),
		bytes.NewBuffer(bodyFiltered),
	)

	for headerKey, headerValue := range dsn.RequestHeaders() {
		request.Header.Set(headerKey, headerValue)
	}

	response, err := i.clientHTTP.Do(request)
	if err != nil {
		log.Warn(fmt.Sprintf("Error client.Do: %v", err))
	}

	log.InfoWithFields("Response received from Sentry", func(e onelog.Entry) {
		e.Int("status", response.StatusCode)
		e.String("forwardTo", request.URL.String())
	})

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warn(fmt.Sprintf("Error reading body: %v", err))
		responseBody = []byte("")
	}

	ctx.SetStatusCode(response.StatusCode)
	ctx.SetBody(responseBody)
}
