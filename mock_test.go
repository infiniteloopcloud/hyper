package hyper

import (
	"context"
	"net/http"
	"testing"
)

func TestMockReadCloser(t *testing.T) {
	body := []byte(`{"id": "test"}`)
	r := NewMockReadCloser(context.Background(), body)

	req, err := http.NewRequest(http.MethodPost, "https://test.example.com", r)
	if err != nil {
		t.Fatal(err)
	}
	defer req.Body.Close()
}
