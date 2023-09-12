package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpConfig struct {
	Headers                map[string]string
	MaskRequestAttributes  []string
	MaskResponseAttributes []string
}

var (
	Delete  func(ctx context.Context, url string, config HttpConfig, response interface{}) error
	Get     func(ctx context.Context, url string, config HttpConfig, response interface{}) error
	Options func(ctx context.Context, url string, config HttpConfig, response interface{}) error
	Post    func(ctx context.Context, url string, data Map, config HttpConfig, response interface{}) error
	Put     func(ctx context.Context, url string, data Map, config HttpConfig, response interface{}) error
)

func init() {
	Delete = delete
	Get = get
	Options = options
	Post = post
	Put = put
}

func delete(ctx context.Context, url string, config HttpConfig, response interface{}) error {
	return httpCall(ctx, http.MethodDelete, url, nil, config, response)
}

func get(ctx context.Context, url string, config HttpConfig, response interface{}) error {
	return httpCall(ctx, http.MethodGet, url, nil, config, response)
}

func options(ctx context.Context, url string, config HttpConfig, response interface{}) error {
	return httpCall(ctx, http.MethodOptions, url, nil, config, response)
}

func post(ctx context.Context, url string, data Map, config HttpConfig, response interface{}) error {
	return httpCall(ctx, http.MethodPost, url, data, config, response)
}

func put(ctx context.Context, url string, data Map, config HttpConfig, response interface{}) error {
	return httpCall(ctx, http.MethodPut, url, data, config, response)
}

func logRequest(ctx context.Context, method string, url string, data Map, config HttpConfig) {
	utils.Debugmf(ctx, Map{
		"eventName": "http-request",
		"method":    method,
		"url":       url,
		"data":      utils.MaskMapKeys(data, config.MaskRequestAttributes...),
		"config": Map{
			"headers": utils.MaskMapKeys(utils.MapStringToMapInterface(config.Headers), "Authorization"),
		},
	}, "Http request to %s %s", method, url)
}

func logResponse(ctx context.Context, method string, url string, data Map, config HttpConfig, responseBytes *[]byte, startTime time.Time) {
	var response map[string]interface{}
	_ = json.Unmarshal(*responseBytes, &response)
	duration := time.Since(startTime)
	utils.Debugmf(ctx, Map{
		"eventName": "http-response",
		"method":    method,
		"url":       url,
		"data":      utils.MaskMapKeys(data, config.MaskRequestAttributes...),
		"config": Map{
			"headers": utils.MaskMapKeys(utils.MapStringToMapInterface(config.Headers), "Authorization"),
		},
		"duration": duration.Milliseconds(),
		"response": utils.MaskMapKeys(response, config.MaskResponseAttributes...),
	}, "Http request to %s %s took %d ms", method, url, duration.Milliseconds())
}

func httpCall(ctx context.Context, method string, url string, data Map, config HttpConfig, response interface{}) error {
	logRequest(ctx, method, url, data, config)
	startTime := time.Now()
	var responseBytes []byte
	defer logResponse(ctx, method, url, data, config, &responseBytes, startTime)
	var req *http.Request
	var err error
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Http_Error, "method": strings.ToLower(method), "url": url, "error": fmt.Sprintf("%v", err)}, "Error creating new http post request - reason: %s", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if config.Headers != nil {
		for k, v := range config.Headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Http_Error, "method": strings.ToLower(method), "url": url, "error": fmt.Sprintf("%v", err)}, "Error requesting %s - reason: %v", url, err)
		return err
	}
	defer res.Body.Close()
	responseBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Http_Error, "method": strings.ToLower(method), "url": url, "error": fmt.Sprintf("%v", err)}, "Fatal error reading response body - reason: %v", err)
		return err
	}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Http_Error, "method": strings.ToLower(method), "url": url, "error": fmt.Sprintf("%v", err)}, "Error decoding response %s - reason: %v", url, err)
		return err
	}
	return nil
}
