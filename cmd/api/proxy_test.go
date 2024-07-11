package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRequests_NotFound(t *testing.T) {
	cmsProxy = newReverseProxy("http://localhost:8081")
	compilerProxy = newReverseProxy("http://localhost:8082")

	cmsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CMS response"))
	}))
	defer cmsServer.Close()

	compilerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Compiler response"))
	}))
	defer compilerServer.Close()

	cmsProxy = newReverseProxy(cmsServer.URL)
	compilerProxy = newReverseProxy(compilerServer.URL)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()

	handleRequests(rr, req)

	wantStatus := http.StatusNotFound
	if status := rr.Code; status != wantStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, wantStatus)
	}

	wantBody := "404 page not found\n"
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(wantBody) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), wantBody)
	}
}
