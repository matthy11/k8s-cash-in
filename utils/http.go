package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"heypay-cash-in-server/constants"
	. "heypay-cash-in-server/utils/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func HttpRespond(ctx context.Context, r *http.Request, w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response, err := json.Marshal(data)
	if err != nil {
		Errorf(nil, "Error encoding data to json - reason: %v", err)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		Errorf(nil, "Error writing data to ResponseWriter - reason: %v", err)
		return
	}
	HttpLogResponse(ctx, r, status, response)
}

func HttpRespondError(ctx context.Context, r *http.Request, w http.ResponseWriter, status int, message, code string) {
	type ErrorResponse struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Code    string `json:"code"`
	}
	response := ErrorResponse{Status: status, Message: message, Code: code}
	HttpRespond(ctx, r, w, status, response)
}

func HttpLogRequest(r *http.Request) {
	ctx := r.Context()
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var data map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &data)
	Infomf(ctx, Map{
		"data":      fmt.Sprintf("%+v", MaskMapKeys(data, "privateKey", "clientSecret")),
		"eventName": "request",
		"method":    r.Method,
		"url":       r.RequestURI,
	}, "Request: %s %s", r.Method, r.RequestURI)
}

func HttpLogResponse(ctx context.Context, r *http.Request, code int, bodyBytes []byte) {
	startTime := ctx.Value(constants.ContextKeyStartTime).(int64)
	now := time.Now().UnixNano()
	elapsedTime := (now - startTime) / 1000000
	var data map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &data)
	success := true
	if code >= http.StatusBadRequest {
		success = false
	}
	LogLastStep(ctx, now)
	Infomf(ctx, Map{
		"code":        code,
		"data":        fmt.Sprintf("%+v", MaskMapKeys(data, "accessToken", "refreshToken")),
		"elapsedTime": elapsedTime,
		"eventName":   "response",
		"method":      r.Method,
		"success":     success,
		"url":         r.RequestURI,
	}, "Response: %s %s", r.Method, r.RequestURI)
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		if strings.Contains(forwarded, ",") {
			return forwarded[:strings.Index(forwarded, ",")]
		}
		return forwarded
	}
	return r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
}
