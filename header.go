package hyper

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/infiniteloopcloud/log"
)

const HeaderPrefix = "Ctx_"

func FromHeader(ctx context.Context, header http.Header, contextKeys map[string]log.ContextField) context.Context {
	for k, h := range header {
		if strings.Contains(k, HeaderPrefix) {
			if len(h) > 0 {
				key := strings.ReplaceAll(k, HeaderPrefix, "")
				if ctx.Value(key) == nil {
					ctx = context.WithValue(ctx, contextKeys[key], h[0])
				}
			}
		}
	}
	return ctx
}

func IntoHeader(ctx context.Context, header http.Header, contextKeys map[string]log.ContextField) http.Header {
	h := copyHeader(header)
	kv := GetValues(ctx, contextKeys)

	for k := range kv {
		if key, ok := k.(log.ContextField); ok {
			h = ElementIntoHeader(ctx, h, key)
		}
	}

	return h
}

func ElementIntoHeader(ctx context.Context, header http.Header, key log.ContextField) http.Header {
	h := copyHeader(header)
	v := ctx.Value(key)
	if value, ok := v.(fmt.Stringer); ok {
		h.Add(HeaderPrefix+string(key), value.String())
	} else if value, ok := v.(string); ok {
		h.Add(HeaderPrefix+string(key), value)
	}
	return h
}

func copyHeader(m http.Header) http.Header {
	cp := make(map[string][]string)
	for k := range m {
		var s = make([]string, len(m[k]))
		copy(s, m[k])
		cp[k] = s
	}

	return cp
}

var castMutex = &sync.Mutex{}

func GetValues(ctx context.Context, contextKeys map[string]log.ContextField) map[interface{}]interface{} {
	castMutex.Lock()
	defer castMutex.Unlock()
	m := make(map[interface{}]interface{})
	for _, field := range contextKeys {
		m[field] = ctx.Value(field)
	}
	// Temporary solution
	// getValuesRecursive(ctx, m)
	return m
}

// type iface struct {
//	itab, data unsafe.Pointer
// }
//
// type valueCtx struct {
//	context.Context
//	key, val interface{}
// }
//
// func getValuesRecursive(ctx context.Context, m map[interface{}]interface{}) {
//	if ctx == nil {
//		return
//	}
//	rv := reflect.ValueOf(ctx)
//	if rv.IsNil() || rv.IsZero() {
//		return
//	}
//	ictxPtr := (*iface)(unsafe.Pointer(&ctx))
//	if ictxPtr == nil {
//		return
//	}
//	ictx := *ictxPtr
//
//	rvUintptr := reflect.ValueOf(ictx.data)
//	if rvUintptr.IsNil() {
//		return
//	}
//
//	valCtx := (*valueCtx)(ictx.data)
//	if valCtx != nil {
//		copyValCtx := *valCtx
//		if copyValCtx.key != nil && copyValCtx.val != nil {
//			if copyValCtx.Context == nil {
//				return
//			}
//			m[copyValCtx.key] = copyValCtx.val
//		}
//
//		if copyValCtx.Context == nil {
//			return
//		}
//		rvValCtxContext := reflect.ValueOf(copyValCtx.Context)
//		if rvValCtxContext.IsNil() || rvValCtxContext.IsZero() {
//			return
//		}
//		getValuesRecursive(copyValCtx.Context, m)
//	}
// }
