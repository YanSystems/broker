package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newReverseProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

var (
	cmsProxy      = newReverseProxy("http://cms:8080")
	compilerProxy = newReverseProxy("http://compiler:8080")
)

func handleRequests(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/cms" || r.URL.Path == "/cms/" || len(r.URL.Path) > 5 && r.URL.Path[:5] == "/cms/":
		cmsProxy.ServeHTTP(w, r)
	case r.URL.Path == "/compiler" || r.URL.Path == "/compiler/" || len(r.URL.Path) > 10 && r.URL.Path[:10] == "/compiler/":
		compilerProxy.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
