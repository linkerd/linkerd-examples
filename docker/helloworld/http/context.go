package http

import (
	"net/http"
	"strings"
)

type linkerdContext map[string]string

func getContext(req *http.Request) *linkerdContext {
	ctx := make(linkerdContext)
	for key, _ := range req.Header {
		if strings.HasPrefix(strings.ToLower(key), "l5d-ctx") {
			ctx[key] = req.Header.Get(key)
		}
	}
	return &ctx
}

func (lc *linkerdContext) withContext(req *http.Request) *http.Request {
	req2 := new(http.Request)
	*req2 = *req
	for key, val := range *lc {
		req2.Header.Set(key, val)
	}
	return req2
}
