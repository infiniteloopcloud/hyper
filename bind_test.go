package hyper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

var (
	ErrMissingRequestBody = errors.New("missing request body")
)

type TestRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (tr TestRequest) Compare(t *testing.T, req TestRequest) {
	if tr.Name != req.Name {
		t.Errorf("Request name should be %s, instead of %s", tr.Name, req.Name)
	}
	if tr.Age != req.Age {
		t.Errorf("Request age should be %d, instead of %d", tr.Age, req.Age)
	}
}

func TestBind(t *testing.T) {
	var scenarios = []struct {
		name        string
		in          string
		expected    TestRequest
		errExpected bool
		err         error
	}{
		{
			name: "fully filled request",
			in:   `{"name":"John","age":44}`,
			expected: TestRequest{
				Name: "John",
				Age:  44,
			},
			errExpected: false,
		},
		{
			name: "missing number field",
			in:   `{"name":"John"}`,
			expected: TestRequest{
				Name: "John",
				Age:  0,
			},
			errExpected: false,
		},
		{
			name: "empty object",
			in:   `{}`,
			expected: TestRequest{
				Name: "",
				Age:  0,
			},
			errExpected: false,
		},
		{
			name:        "null",
			in:          `null`,
			errExpected: false,
		},
		{
			name:        "missing body",
			in:          "",
			errExpected: true,
			err:         ErrMissingRequestBody,
		},
	}
	ctx := context.Background()
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			var reqBody io.Reader
			if scenario.in != "" {
				reqBody = bytes.NewBufferString(scenario.in)
			}
			req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:9999", reqBody)
			if err != nil {
				t.Fatal(err)
			}
			var body TestRequest
			if err := Bind(ctx, req, &body); (err != nil) != scenario.errExpected && !errors.Is(err, scenario.err) {
				t.Errorf("Expected error is %v but we got %v", scenario.errExpected, err)
			}
			scenario.expected.Compare(t, body)
		})
	}
}
