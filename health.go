package hyper

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/infiniteloopcloud/log"
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

func ReadyzWithChecks(checks ...func(ctx context.Context) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		for _, check := range checks {
			if err := check(r.Context()); err != nil {
				log.Error(r.Context(), err, "check in readyz")
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Hello returns http.StatusOk and the Hello response
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	w.Write([]byte(`{"msg":"Hello"}`))
}

func ServiceLive(ctx context.Context, url string) error {
	resp, err := Request(ctx, nil, RequestOpts{
		BaseURL:  url,
		Method:   http.MethodGet,
		Endpoint: "/livez",
		Client:   Client(),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("invalid status code %d", resp.StatusCode)
		log.Error(ctx, err, "service not alive error")
		return err
	}
	return nil
}

func ServiceReady(ctx context.Context, url string) error {
	resp, err := Request(ctx, nil, RequestOpts{
		BaseURL:  url,
		Method:   http.MethodGet,
		Endpoint: "/readyz",
		Client:   Client(),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("invalid status code %d", resp.StatusCode)
		log.Error(ctx, err, "service readiness error")
		return err
	}
	return nil
}
