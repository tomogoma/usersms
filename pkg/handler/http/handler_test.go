package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/logging"
	testingH "github.com/tomogoma/usersms/pkg/mocks"
)

func TestNewHandler(t *testing.T) {
	tt := []struct {
		name           string
		guard          Guard
		logger         logging.Logger
		allowedOrigins []string
		expErr         bool
	}{
		{
			name:           "valid deps",
			guard:          &testingH.Guard{},
			logger:         &testingH.Logger{},
			allowedOrigins: []string{"*"},
			expErr:         false,
		},
		{
			name:   "valid deps (nil origins)",
			guard:  &testingH.Guard{},
			logger: &testingH.Logger{},
			expErr: false,
		},
		{
			name:   "nil guard",
			guard:  nil,
			logger: &testingH.Logger{},
			expErr: true,
		},
		{
			name:   "nil logger",
			guard:  &testingH.Guard{},
			logger: nil,
			expErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h, err := NewHandler(tc.guard, tc.logger, "", tc.allowedOrigins)
			if tc.expErr {
				if err == nil {
					t.Fatal("Expected an error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if h == nil {
				t.Fatalf("http.NewHandler() yielded a nil handler!")
			}
		})
	}
}

func TestHandler_handleRoute(t *testing.T) {
	tt := []struct {
		name          string
		baseURL       string
		reqURLSuffix  string
		reqMethod     string
		reqBody       string
		reqWBasicAuth bool
		expStatusCode int
		guard         Guard
	}{
		// values starting and ending with "_" are place holders for variables
		// e.g. _loginType_ is a place holder for "any (valid) login type"

		{
			name:          "status",
			guard:         &testingH.Guard{},
			reqURLSuffix:  "/status",
			reqMethod:     http.MethodGet,
			expStatusCode: http.StatusOK,
		},
		{
			name:          "status guard error",
			guard:         &testingH.Guard{ExpAPIKValidErr: errors.Newf("guard error")},
			reqURLSuffix:  "/status",
			reqMethod:     http.MethodGet,
			expStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			lg := &testingH.Logger{}
			h := newHandler(t, tc.guard, lg, tc.baseURL, nil)
			srvr := httptest.NewServer(h)
			defer srvr.Close()

			req, err := http.NewRequest(
				tc.reqMethod,
				srvr.URL+tc.reqURLSuffix,
				bytes.NewReader([]byte(tc.reqBody)),
			)
			if err != nil {
				t.Fatalf("Error setting up: new request: %v", err)
			}
			if tc.reqWBasicAuth {
				req.SetBasicAuth("username", "password")
			}

			cl := &http.Client{}
			resp, err := cl.Do(req)
			if err != nil {
				lg.PrintLogs(t)
				t.Fatalf("Do request error: %v", err)
			}

			if resp.StatusCode != tc.expStatusCode {
				lg.PrintLogs(t)
				t.Errorf("Expected status code %d, got %s",
					tc.expStatusCode, resp.Status)
			}
		})
	}
}

func newHandler(t *testing.T, g Guard, lg logging.Logger, baseURL string, allowedOrigins []string) http.Handler {
	h, err := NewHandler(g, lg, baseURL, allowedOrigins)
	if err != nil {
		t.Fatalf("http.NewHandler(): %v", err)
	}
	return h
}
