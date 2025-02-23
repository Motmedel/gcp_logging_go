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

func ParseHttp(request *http.Request, response *http.Response) *LogEntry {
	if request == nil && response == nil {
		return nil
	}

	var httpRequest HttpRequest

	if request != nil {
		httpRequest.RequestMethod = request.Method
		httpRequest.UserAgent = request.UserAgent()
		httpRequest.RemoteIp = request.RemoteAddr
		httpRequest.Referer = request.Referer()
		httpRequest.Protocol = fmt.Sprintf("HTTP/%d.%d", request.ProtoMajor, request.ProtoMinor)
	}

	if response != nil {
		httpRequest.Status = response.StatusCode
	}

	return &LogEntry{HttpRequest: &httpRequest}
}
