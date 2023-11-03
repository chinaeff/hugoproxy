package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReverseProxyMiddleware(t *testing.T) {
	fakeHugo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fake Hugo Response"))
	}))
	defer fakeHugo.Close()

	_, port, _ := net.SplitHostPort(strings.TrimPrefix(fakeHugo.URL, "http://"))
	host := "localhost"

	proxy := NewReverseProxy(host, port)

	req := httptest.NewRequest("GET", "/api/", nil)
	w := httptest.NewRecorder()

	proxy.ReverseProxy(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach this point")
	})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	expectedBody := "Hello from API"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body: %s, got %s", expectedBody, w.Body.String())
	}

	req = httptest.NewRequest("GET", "/other/path", nil)
	w = httptest.NewRecorder()

	proxy.ReverseProxy(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach this point")
	})).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	expectedBody = "Fake Hugo Response"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body: %s, got %s", expectedBody, w.Body.String())
	}
}

func TestWorkerTest(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		WorkerTest()
	}()

	time.Sleep(5 * time.Second)

	close(done)
}

func TestNewReverseProxy(t *testing.T) {
	host := "example.com"
	port := "8080"

	proxy := NewReverseProxy(host, port)

	if proxy.host != host {
		t.Errorf("Expected host to be %s, but got %s", host, proxy.host)
	}

	if proxy.port != port {
		t.Errorf("Expected port to be %s, but got %s", port, proxy.port)
	}
}
