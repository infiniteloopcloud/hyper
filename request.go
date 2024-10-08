package hyper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/infiniteloopcloud/log"
)

func checkBaseURL(baseURL string) error {
	if baseURL == "" {
		return errors.New("base url is empty, set it with SetBaseURL(...)")
	}

	return nil
}

func endpoint(baseURL, s string) string {
	return baseURL + s
}

type RequestOpts struct {
	BaseURL       string
	Method        string
	Endpoint      string
	Request       interface{}
	Headers       map[string]string
	ContextKeys   map[string]log.ContextField
	SkipBodyClose bool

	Client *http.Client
}

//nolint:gocritic
func Request(ctx context.Context, respStruct interface{}, opts RequestOpts) (*http.Response, error) {
	if err := checkBaseURL(opts.BaseURL); err != nil {
		return nil, err
	}

	var r *http.Request
	var err error
	if req, ok := opts.Request.(io.Reader); ok {
		r, err = http.NewRequest(opts.Method, endpoint(opts.BaseURL, opts.Endpoint), req)
	} else if opts.Request != nil {
		b := new(bytes.Buffer)
		if err := jsonEncoder.Encode(b, opts.Request); err != nil {
			return nil, err
		}
		r, err = http.NewRequest(opts.Method, endpoint(opts.BaseURL, opts.Endpoint), b)
	} else {
		r, err = http.NewRequest(opts.Method, endpoint(opts.BaseURL, opts.Endpoint), nil)
	}
	if err != nil {
		return nil, err
	}
	r.Header = IntoHeader(ctx, r.Header, opts.ContextKeys)

	if opts.Headers != nil {
		for k, v := range opts.Headers {
			r.Header.Add(k, v)
		}
	}
	resp, err := opts.Client.Do(r)
	if err != nil {
		return nil, err
	}

	if !opts.SkipBodyClose {
		defer resp.Body.Close()
	}

	if respStruct != nil {
		if err := jsonEncoder.Decode(resp.Body, respStruct); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func SilentProxy(ctx context.Context, w http.ResponseWriter, resp *http.Response) {
	for key, valueSlice := range resp.Header {
		for _, v := range valueSlice {
			w.Header().Add(key, v)
		}
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		defer resp.Body.Close()
		return
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		log.Error(ctx, err, "error copying resp to buffer")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	if _, err := w.Write(b.Bytes()); err != nil {
		log.Error(ctx, err, "error copying resp to buffer")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
