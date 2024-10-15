package hello_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	service "github.com/go-obvious/server-example/internal/service/hello"
	"github.com/go-obvious/server/test"
)

func TestHelloServiceEndpoint(t *testing.T) {

	tcs := []struct {
		Description string

		// Request / setup
		URL           string
		RequestHeader map[string]string

		// Response
		Code   int
		Header map[string]string
		Method string
		Body   string // request body to send
		Want   string // response body to expect
	}{
		{
			Description: "lookup returns 200 dummy test",
			URL:         "/hello",
			Method:      http.MethodGet,
			Code:        http.StatusOK,
			Header:      map[string]string{},
		},
	}

	for _, tc := range tcs {
		testf := func(t *testing.T) {
			req := http.Request{
				Method: tc.Method,
				Body:   io.NopCloser(strings.NewReader(tc.Body)),
				Header: map[string][]string{},
			}

			//
			// set any headers
			for k, v := range tc.RequestHeader {
				req.Header.Set(k, v)
			}

			//
			// Log the HTTP call we are about to test
			t.Log(req.Method, req.URL)

			resp, err := test.InvokeService(
				service.NewService("/hello").Service,
				"/hello",
				req,
			)
			assert.NoError(t, err)

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Reading response body: %v", err)
			}
			if resp.StatusCode != tc.Code {
				t.Errorf("Incorrect status code, got %d, want %d; body: %s", resp.StatusCode, tc.Code, body)
			}

			for k, v := range tc.Header {
				r := resp.Header.Get(k)
				if r != v {
					t.Errorf("Incorrect header %q received, got %q, want %q", k, r, v)
				}
			}

			if tc.Want != "" && string(body) != tc.Want {
				t.Errorf("Incorrect response body, got %q, want %q", body, tc.Want)
			}
		}
		t.Run(tc.Description, testf)
	}
}
