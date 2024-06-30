package gcp_logging

import (
	"fmt"
	"net/http"
	"time"
)

func NewDuration(duration *time.Duration) *Duration {
	return &Duration{
		Seconds: int(*duration / time.Second),
		Nanos:   int(*duration % time.Second),
	}
}

func ParseHttpRequest(request *http.Request) (*HttpRequest, error) {
	return &HttpRequest{
		RequestMethod: request.Method,
		UserAgent:     request.UserAgent(),
		RemoteIp:      request.RemoteAddr,
		Referer:       request.Referer(),
		Protocol:      fmt.Sprintf("HTTP/%d.%d", request.ProtoMajor, request.ProtoMinor),
	}, nil
}
