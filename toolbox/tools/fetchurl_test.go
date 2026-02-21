package tools

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// patchDefaultTransport replaces http.DefaultTransport for the duration of the
// test and restores it afterwards. FetchURL creates an http.Client with a nil
// Transport, which falls back to http.DefaultTransport at runtime.
func patchDefaultTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()
	orig := http.DefaultTransport
	http.DefaultTransport = rt
	t.Cleanup(func() { http.DefaultTransport = orig })
}

func TestFetchURL_ValidationErrors(t *testing.T) {
	cases := []struct {
		name         string
		params       Parameters
		wantErrMatch string
	}{
		{
			name:         "missing url param",
			params:       Parameters{},
			wantErrMatch: "missing required parameter: url",
		},
		{
			name:         "empty url param",
			params:       Parameters{"url": ""},
			wantErrMatch: "missing required parameter: url",
		},
		{
			name:         "whitespace url param",
			params:       Parameters{"url": "   "},
			wantErrMatch: "missing required parameter: url",
		},
		{
			name:         "http url rejected",
			params:       Parameters{"url": "http://example.com"},
			wantErrMatch: "only HTTPS URLs are allowed",
		},
		{
			name:         "non-http url rejected",
			params:       Parameters{"url": "ftp://example.com"},
			wantErrMatch: "only HTTPS URLs are allowed",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := FetchURL(context.Background(), testCase.params)
			assert.Error(t, err)
			assert.ErrorContains(t, err, testCase.wantErrMatch)
		})
	}
}

func TestFetchURL_Success(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "hello from test server")
	}))
	t.Cleanup(server.Close)
	patchDefaultTransport(t, server.Client().Transport)

	result, err := FetchURL(context.Background(), Parameters{"url": server.URL})

	// The test server URL starts with https:// since it is a TLS server.
	assert.NoError(t, err)
	assert.Equal(t, "hello from test server", result)
}

func TestFetchURL_NonSuccessStatusCode(t *testing.T) {
	cases := []struct {
		name   string
		status int
	}{
		{name: "404 Not Found", status: http.StatusNotFound},
		{name: "500 Internal Server Error", status: http.StatusInternalServerError},
		{name: "403 Forbidden", status: http.StatusForbidden},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(testCase.status)
			}))
			t.Cleanup(server.Close)
			patchDefaultTransport(t, server.Client().Transport)

			_, err := FetchURL(context.Background(), Parameters{"url": server.URL})
			assert.Error(t, err)
			assert.ErrorContains(t, err, fmt.Sprintf("status code %d", testCase.status))
		})
	}
}

func TestFetchURL_ResponseTruncatedAtMaxOutputSize(t *testing.T) {
	body := strings.Repeat("a", MaxOutputSize+100)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	patchDefaultTransport(t, server.Client().Transport)

	result, err := FetchURL(context.Background(), Parameters{"url": server.URL})
	assert.NoError(t, err)
	assert.Len(t, result, MaxOutputSize)
}

func TestFetchURL_ContextCancelled(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	t.Cleanup(server.Close)
	patchDefaultTransport(t, server.Client().Transport)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchURL(ctx, Parameters{"url": server.URL})
	assert.Error(t, err)
}
