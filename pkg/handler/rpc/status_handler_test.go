package rpc_test

import (
	"testing"

	"github.com/tomogoma/usersms/pkg/handler/rpc"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/mocks"
	"context"
	"github.com/tomogoma/usersms/pkg/api"
	"github.com/tomogoma/go-typed-errors"
)

func TestNewHandler(t *testing.T) {
	tt := []struct {
		name   string
		guard  rpc.Guard
		logger logging.Logger
		expErr bool
	}{
		{
			name:   "valid deps",
			guard:  &mocks.Guard{},
			logger: &mocks.Logger{},
			expErr: false,
		},
		{
			name:   "nil guard",
			guard:  nil,
			logger: &mocks.Logger{},
			expErr: true,
		},
		{
			name:   "nil logger",
			guard:  &mocks.Guard{},
			logger: nil,
			expErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sh, err := rpc.NewStatusHandler(tc.guard, tc.logger)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if sh == nil {
				t.Fatalf("Got nil *rpc.StatusHandler")
			}
		})
	}
}


func TestStatusHandler_Check(t *testing.T) {
	tt := []struct {
		name string
		guard *mocks.Guard
		req *api.Request
		expErr bool
	}{
		{
			name: "valid",
			guard: &mocks.Guard{},
			req: &api.Request{},
			expErr: false,
		},
		{
			name: "forbidden",
			guard: &mocks.Guard{ExpAPIKValidErr: errors.NewForbidden("guard")},
			req: &api.Request{},
			expErr: true,
		},
		{
			name: "unauthorized",
			guard: &mocks.Guard{ExpAPIKValidErr: errors.NewUnauthorized("guard")},
			req: &api.Request{},
			expErr: true,
		},
		{
			name: "internal error",
			guard: &mocks.Guard{ExpAPIKValidErr: errors.Newf("guard")},
			req: &api.Request{},
			expErr: true,
		},
	}
	for _, tc := range tt {
		sh := newStatusHandler(t, tc.guard, &mocks.Logger{})
		resp := new(api.Response)
		err := sh.Check(context.TODO(), tc.req, resp)
		if tc.expErr {
			if err ==nil {
				t.Fatalf("Expected an error, got nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
	}
}

func newStatusHandler(t *testing.T, g rpc.Guard, lg logging.Logger) *rpc.StatusHandler {
	sh, err := rpc.NewStatusHandler(g, lg)
	if err != nil {
		t.Fatalf("Error setting up: new status handler: %v", err)
	}
	return sh
}