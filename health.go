package hyper

import (
	"net/http"
	"sync/atomic"
)

// isReady defines whether the api is ready for request
var isReady *atomic.Value

func init() {
	isReady = &atomic.Value{}
}

// Ready enables the server to accept connections
func Ready() {
	isReady.Store(true)
}

// NotReady disables the server to accept connections
func NotReady() {
	isReady.Store(false)
}

// Livez returns http.StatusOk when the service does not need to be killed
func Livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz returns http.StatusOk when the service is ready accept connections, otherwise Http.StatusServiceUnavailable
func Readyz(w http.ResponseWriter, r *http.Request) {
	if !isReady.Load().(bool) {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Hello returns http.StatusOk and the Hello response
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	w.Write([]byte(`{"msg":"Hello"}`))
}
