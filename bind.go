package hyper

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/infiniteloopcloud/log"
	"net/http"
)

var (
	ErrBindMissingBody = errors.New("request bind error, missing request body")
	ErrBind            = "request bind error, "
)

var (
	MissingRequestBody = "Missing request body"
	InvalidRequest     = "Invalid request"
)

const (
	CorrelationID log.ContextField = "correlation_id"
)

// Bind is binding the request body into v variable
func Bind(ctx context.Context, r *http.Request, v interface{}) error {
	if r.Body == nil {
		return BadRequest(ctx, MissingRequestBody, ErrBindMissingBody)
	}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return BadRequest(ctx, InvalidRequest, errors.New(ErrBind+err.Error()))
	}
	return nil
}

func GetCorrelationID(ctx context.Context) string {
	ctxVal := ctx.Value(CorrelationID)
	if v, ok := ctxVal.(string); ok {
		return v
	}
	return ""
}
