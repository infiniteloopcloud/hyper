package hyper

import (
	"bytes"
	"context"
	"net/http"

	"github.com/infiniteloopcloud/go/weird"
	"github.com/infiniteloopcloud/log"
)

const headerContentType = "Content-Type"

// Writer wraps the http.ResponseWriter with status code
// This wrapper is necessary because we want to store the statusCode,
// (needed for the prometheus 'responseStatus' histogram collector)
// and response body (needed for logging middleware)

// Writer also implements the http.ResponseWriter interface
// In the future we may add other fields as well.
type Writer struct {
	W          http.ResponseWriter
	body       *bytes.Buffer
	StatusCode int
}

// NewWriter creates a new Writes
func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		W:    w,
		body: &bytes.Buffer{},
	}
}

// Header returns the header map that will be sent by WriteHeader
func (r *Writer) Header() http.Header {
	return r.W.Header()
}

// Write writes the data to the connection as part of an HTTP reply, and stores the
// response data for logging purpose
func (r *Writer) Write(bytes []byte) (int, error) {
	r.body.Write(bytes)
	return r.W.Write(bytes)
}

// WriteHeader 'overrides' the http.ResponseWriter's WriteHeader method
func (r *Writer) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.W.WriteHeader(statusCode)
}

func (r *Writer) ResponseBody() string {
	return r.body.String()
}

// Wrapper wraps the response
type Wrapper struct {
	// TokenJWT returns a new JWT token, if it filled that means the frontend should change it
	// NOTE: temporary it's hidden, because we use cookies
	TokenJWT string `json:"-"`

	// UserJWT returns a new JWT with user info, if it filled that means the frontend should change it
	// NOTE: temporary it's hidden, because we use cookies
	UserJWT string `json:"-"`

	// Error
	Error string `json:"error,omitempty"`
	// Data is the exact data of the response
	Data interface{} `json:"data,omitempty"`
}

// Success build a success response with HTTP Status OK
func Success(ctx context.Context, w http.ResponseWriter, val interface{}) {
	Generic(ctx, w, val, http.StatusOK)
}

// SuccessDownload build a success response with content type octet/stream
func SuccessDownload(ctx context.Context, w http.ResponseWriter, res []byte) {
	w.Header().Set(headerContentType, "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment;`)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		log.Error(ctx, err, "unable to write response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Created build a success response with HTTP Status Created
func Created(ctx context.Context, w http.ResponseWriter, val interface{}) {
	Generic(ctx, w, val, http.StatusCreated)
}

// NoContent build a no content success response
func NoContent(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error build an error response, the HTTP status will get from the error, default 500
func Error(ctx context.Context, w http.ResponseWriter, err error) {
	if e, ok := err.(weird.Error); ok {
		Generic(ctx, w, e, e.StatusCode)
	} else {
		Error(ctx, w, internalServerError(ctx, err))
	}
}

// Generic build a generic response
func Generic(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) {
	if ok := writeResponseHeaderAndBody(ctx, w, val, statusCode); !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// writeResponseHeaderAndBody wrapping the response data into the Wrapper instance
// return true if the marshaling and the response writing doesn't occur error
// otherwise false
func writeResponseHeaderAndBody(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) bool {
	var data Wrapper
	if err, ok := val.(error); ok {
		if werr, wok := err.(weird.Error); wok && werr.Msg != "" {
			data.Error = werr.Msg + " (" + GetCorrelationID(ctx) + ")"
		} else {
			data.Error = "Unknown error (" + GetCorrelationID(ctx) + ")"
		}
	} else {
		data.Data = val
	}
	res := new(bytes.Buffer)
	err := jsonEncoder.Encode(res, data)
	if err != nil {
		log.Error(ctx, err, "unable to marshal response")
		return false
	}

	w.Header().Set(headerContentType, "application/json")

	// if data.TokenJWT != "" && data.UserJWT != "" {
	//	cookie.Login(w, data.TokenJWT, data.UserJWT)
	// } // TODO move this out

	w.WriteHeader(statusCode)
	if resBytes := res.Bytes(); resBytes != nil {
		_, err = w.Write(resBytes)
		if err != nil {
			log.Error(ctx, err, "unable to write response")
			return false
		}
	}

	return true
}

func internalServerError(_ context.Context, err error) error {
	if weirdErr, ok := err.(weird.Error); ok {
		weirdErr.InnerError = err
		return weirdErr
	}
	return weird.New("", err, http.StatusInternalServerError)
}
