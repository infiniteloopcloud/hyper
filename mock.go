package hyper

import (
	"bytes"
	"context"

	"github.com/infiniteloopcloud/log"
)

// MockReadCloser implements the io.ReadCloser interface
// Mocks the request.Body
type MockReadCloser struct {
	*bytes.Buffer
}

func NewMockReadCloser(ctx context.Context, data interface{}) MockReadCloser {
	buffer := bytes.Buffer{}
	if data != nil {
		b := new(bytes.Buffer)
		err := jsonEncoder.Encode(b, data)
		if err != nil {
			log.Error(ctx, err, "error marshaling data")
		}
		buffer.Write(b.Bytes())
	}

	return MockReadCloser{
		Buffer: &buffer,
	}
}

func (cb MockReadCloser) Close() (err error) {
	return nil
}
