package hyper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/infiniteloopcloud/log"
)

var (
	ErrBindMissingBody = errors.New("request bind error, missing request body")
	ErrBind            = "request bind error, "
)

var (
	MissingRequestBody = "Missing request body"
	InvalidRequest     = "Invalid request"
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
	ctxVal := ctx.Value(log.CorrelationID)
	if v, ok := ctxVal.(string); ok {
		return v
	}
	return ""
}
