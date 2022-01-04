package hyper

import (
	"context"
	"errors"
	"net/http"
)

const defaultAddress = ":8000"

var ErrMissingHandler = errors.New("missing http.Handler, pass it in the HTTP struct")

type HTTP struct {
	TLS        TLSDescriptor
	TLSEnabled bool

	Address string
	Handler http.Handler

	Log Logger

	PreHooks  []func(ctx context.Context)
	PostHooks []func(ctx context.Context)
}

type Logger struct {
	Infof  infoLog
	Errorf errorLog
}

type infoLog func(ctx context.Context, format string, args ...interface{})
type errorLog func(ctx context.Context, err error, format string, args ...interface{})

func (h HTTP) Serve(ctx context.Context, finisher chan error) {
	if h.Address == "" {
		h.Address = defaultAddress
	}
	if h.Handler == nil {
		finisher <- ErrMissingHandler
		return
	}

	srv := http.Server{
		Addr:    h.Address,
		Handler: h.Handler,
	}

	if h.TLSEnabled {
		var err error
		srv.TLSConfig, err = h.TLS.from()
		if err != nil {
			finisher <- ErrMissingHandler
			return
		}
	}

	go h.shutdownListener(ctx, &srv, finisher)

	if h.TLSEnabled {
		h.Log.Infof(ctx, "HTTP server listening on %s with TLS", h.Address)
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			finisher <- err
		}
	} else {
		h.Log.Infof(ctx, "HTTP server listening on %s", h.Address)
		if err := srv.ListenAndServe(); err != nil {
			finisher <- err
		}
	}
}

func (h HTTP) shutdownListener(ctx context.Context, srv *http.Server, finisher chan error) {
	<-ctx.Done()
	h.Log.Infof(ctx, "graceful shutdown finished... exiting")
	finisher <- srv.Shutdown(ctx)
}
