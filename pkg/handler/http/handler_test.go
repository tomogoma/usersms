package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/mocks"
)

func TestNewHandler(t *testing.T) {
	tt := []struct {
		name   string
		conf   Config
		expErr bool
	}{
		{
			name: "valid deps",
			conf: Config{
				Guard:        &mocks.Guard{},
				Logger:       &mocks.Logger{},
				Rater:        &mocks.Rater{},
				UserProfiler: &mocks.User{},
			},
			expErr: false,
		},
		{
			name: "valid deps (nil origins)",
			conf: Config{
				Guard:  &mocks.Guard{},
				Logger: &mocks.Logger{},
				Rater:        &mocks.Rater{},
				UserProfiler: &mocks.User{},
			},
			expErr: false,
		},
		{
			name: "nil guard",
			conf: Config{
				Guard:  nil,
				Logger: &mocks.Logger{},
				Rater:        &mocks.Rater{},
				UserProfiler: &mocks.User{},
			},
			expErr: true,
		},
		{
			name: "nil logger",
			conf: Config{
				Guard:  &mocks.Guard{},
				Logger: nil,
				Rater:        &mocks.Rater{},
				UserProfiler: &mocks.User{},
			},
			expErr: true,
		},
		{
			name: "nil Rater",
			conf: Config{
				Guard:        &mocks.Guard{},
				Logger:       &mocks.Logger{},
				Rater:        nil,
				UserProfiler: &mocks.User{},
			},
			expErr: true,
		},
		{
			name: "nil UserProfiler",
			conf: Config{
				Guard:        &mocks.Guard{},
				Logger:       &mocks.Logger{},
				Rater:        &mocks.Rater{},
				UserProfiler: nil,
			},
			expErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h, err := NewHandler(tc.conf)
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
		reqURLSuffix  string
		reqMethod     string
		reqBody       string
		reqWBasicAuth bool
		expStatusCode int
		conf          Config
	}{
		{
			name:          "status",
			conf:          Config{Guard: &mocks.Guard{}},
			reqURLSuffix:  "/status",
			reqMethod:     http.MethodGet,
			expStatusCode: http.StatusOK,
		},
		{
			name: "status guard error",
			conf: Config{
				Guard: &mocks.Guard{ExpAPIKValidErr: errors.Newf("guard error")},
			},
			reqURLSuffix:  "/status",
			reqMethod:     http.MethodGet,
			expStatusCode: http.StatusInternalServerError,
		},
		{
			name:          "not found",
			conf:          Config{Guard: &mocks.Guard{}},
			reqURLSuffix:  "/none_existent",
			reqMethod:     http.MethodGet,
			expStatusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			lg := &mocks.Logger{}
			tc.conf.Logger = lg
			tc.conf.Rater = &mocks.Rater{}
			tc.conf.UserProfiler = &mocks.User{}
			h := newHandler(t, tc.conf)
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

func newHandler(t *testing.T, conf Config) http.Handler {
	h, err := NewHandler(conf)
	if err != nil {
		t.Fatalf("http.NewHandler(): %v", err)
	}
	return h
}
