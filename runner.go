package hyper

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func (h HTTP) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	for _, fn := range h.PreHooks {
		fn(ctx)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	finisher := make(chan error, 1)
	go h.Serve(ctx, finisher)
	go h.osSignalListener(ctx, cancelFunc, sig)

	if err := <-finisher; err != nil {
		h.Log.Errorf(ctx, err, "failed to run api")
	}

	for _, fn := range h.PostHooks {
		fn(ctx)
	}
}

func (h HTTP) osSignalListener(ctx context.Context, cancelFunc func(), sig chan os.Signal) {
	osCall := <-sig
	h.Log.Infof(ctx, "%s received, graceful shutdown started", osCall.String())
	cancelFunc()
}
