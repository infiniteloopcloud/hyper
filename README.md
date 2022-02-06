# Hyper

General HTTP helper library aims to be customizable.

- [Non-TLS](#non-tls)
- [TLS](#tls)
- [Bind request body](#bind-request-body)
- [Error handling](#error-handling)
- [Client](#client)

### Usage

#### Non-TLS

```go
package main

import (
	"context"
	"os"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

func main() {
	var h = hyper.HTTP{
		PreHooks: []func(ctx context.Context){
			preHookSetLogLevel,
		},
		Log: hyper.Logger{
			Infof:  log.Infof,
			Errorf: log.Errorf,
		},
		Address: os.Getenv("SERVICE_HTTP_ADDRESS"),
		Handler: chi.NewRouter(),
	}
	
	h.Run()
}

func preHookSetLogLevel(_ context.Context) {
	log.SetLevel(func() uint8 {
		return 1
	}())
}
```

#### TLS

```go
package main

import (
	"context"
	"os"

	"github.com/infiniteloopcloud/hyper"
	"github.com/infiniteloopcloud/log"
)

func main() {
	var h = hyper.HTTP{
		PreHooks: []func(ctx context.Context){
			preHookSetLogLevel,
		},
		Log: hyper.Logger{
			Infof:  log.Infof,
			Errorf: log.Errorf,
		},
		Address: os.Getenv("SERVICE_HTTP_ADDRESS"),
		Handler: chi.NewRouter(),
	}

	// NewEnvironmentTLS able to build PEM block based on the provided information
	if os.Getenv("SERVICE_TLS_SUPPORT") == "PEM_BUILDER" {
		h.TLSEnabled = true
		h.TLS = hyper.NewEnvironmentTLS(hyper.EnvironmentTLSOpts{
			TLSCert:          "SERVICE_TLS_CERT",
			TLSCertBlockType: "SERVICE_TLS_CERT_BLOCK_NAME",
			TLSKey:           "SERVICE_TLS_KEY",
			TLSKeyBlockType:  "SERVICE_TLS_KEY_BLOCK_NAME",
		})
    // NewFileTLS reads the files from the provided location
	} else if os.Getenv("SERVICE_TLS_SUPPORT") == "FILE" {
		h.TLSEnabled = true
		h.TLS = hyper.NewFileTLS(hyper.FileTLSOpts{
			TLSCertPath: "SERVICE_TLS_CERT",
			TLSKeyPath:  "SERVICE_TLS_KEY",
		})
	}
	
	h.Run()
}

func preHookSetLogLevel(_ context.Context) {
	log.SetLevel(func() uint8 {
		return 1
	}())
}
```

#### Bind request body

```go
package main

import (
	"net/http"

	"github.com/infiniteloopcloud/hyper"
)

type Request struct {
	Field1 string `json:"field1"`
}

func main() {}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqStruct Request
	// Bind the JSON request
	if err := hyper.Bind(ctx, r, &reqStruct); err != nil {
		hyper.ReturnBadRequest(ctx, w, "invalid request body", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
```

#### Error handling

There are pure error wrappers which wraps the error into a weird.Error:

- `BadRequest(...)`
- `NotFound(...)`
- `Unauthorized(...)`
- `Forbidden(...)`
- `InternalServerError(...)`

All of these has a Return{...} function if the function has access to `http.ResponseWriter`

- `ReturnBadRequest(...)`
- `ReturnNotFound(...)`
- `ReturnUnauthorized(...)`
- `ReturnForbidden(...)`
- `ReturnInternalServerError(...)`

```go
package main

import (
	"net/http"

	"github.com/infiniteloopcloud/hyper"
)

type Request struct {
	Field1 string `json:"field1"`
}

func main() {}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqStruct Request
	if err := hyper.Bind(ctx, r, &reqStruct); err != nil {
		// Return bad request will write the response into w
		hyper.ReturnBadRequest(ctx, w, "invalid request body", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
```

#### Client

HTTP client-side helper functions.

```go
package main

import (
	"context"

	"github.com/infiniteloopcloud/hyper"
)

type Response struct {
	Field1 string `json:"f"`
}

func main() {
	ctx := context.Background()
	var resp Response
	hyper.Request(ctx, &resp, hyper.RequestOpts{
		Method:         "POST",
		Endpoint:       "/api/endpoint",
		// RequestStruct:  req,
		Client:         hyper.Client(),
	})
	// handle err
	// use resp
}
```