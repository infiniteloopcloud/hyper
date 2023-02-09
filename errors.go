package hyper

import (
	"context"
	"net/http"

	"github.com/infiniteloopcloud/go/weird"
)

func BadRequest(ctx context.Context, msg string, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New(msg, err, http.StatusBadRequest)
}

func ReturnBadRequest(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	Error(ctx, w, BadRequest(ctx, msg, err))
}

func NotFound(ctx context.Context, msg string, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New(msg, err, http.StatusNotFound)
}

func ReturnNotFound(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	Error(ctx, w, NotFound(ctx, msg, err))
}

func Unauthorized(ctx context.Context, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New("", err, http.StatusUnauthorized)
}

func ReturnUnauthorized(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, Unauthorized(ctx, err))
}

func Forbidden(ctx context.Context, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New("", err, http.StatusForbidden)
}

func ReturnForbidden(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, Forbidden(ctx, err))
}

func InternalServerError(ctx context.Context, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New("", err, http.StatusInternalServerError)
}

func ReturnInternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, InternalServerError(ctx, err))
}
