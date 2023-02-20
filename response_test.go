package hyper

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestSuccess(t *testing.T) {
	var ctx = context.Background()
	var rw = &fakeResponseWriter{}
	var s = struct {
		Field1 string
		Field2 int
	}{
		Field1: "test",
		Field2: 12,
	}
	Success(ctx, rw, s)
	var expected = `{"data":{"Field1":"test","Field2":12}}`
	if !strings.Contains(string(rw.data), expected) {
		t.Errorf("ResponseWriter data should contain %s, instead of %s", expected, string(rw.data))
	}
	if rw.statusCode != 200 {
		t.Errorf("ResponseWriter status code should be 200, instead of %d", rw.statusCode)
	}
}

type fakeResponseWriter struct {
	data       []byte
	statusCode int
}

func (f *fakeResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (f *fakeResponseWriter) Write(bytes []byte) (int, error) {
	f.data = bytes
	return 0, nil
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}
