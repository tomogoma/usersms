package roach_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	apiH "github.com/tomogoma/go-api-guard"
	"github.com/tomogoma/usersms/pkg/api"
	"github.com/tomogoma/usersms/pkg/db/roach"
)

func TestRoach_InsertAPIKey(t *testing.T) {
	setupTime := time.Now()
	conf, tearDown := setup(t)
	defer tearDown()
	r := newRoach(t, conf)
	validKey := []byte(strings.Repeat("axui", 14))
	usrID := "123"
	tt := []struct {
		testName string
		key      []byte
		usrID    string
		expErr   bool
	}{
		{testName: "valid", key: validKey, usrID: usrID, expErr: false},
		{testName: "bad user ID", key: validKey, usrID: "bad id", expErr: true},
		{testName: "empty key", key: []byte{}, usrID: usrID, expErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			retI, err := r.InsertAPIKey(tc.usrID, tc.key)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if retI == nil {
				t.Fatalf("Got nil group")
			}
			ret, ok := retI.(api.Key)
			if !ok {
				t.Fatalf("Exptected API key of type %T, got %T",
					api.Key{}, ret)
			}
			if ret.ID == "" {
				t.Errorf("ID was not assigned")
			}
			if ret.LastUpdated.Before(setupTime) {
				t.Errorf("UpdateDate was not assigned")
			}
			if ret.Created.Before(setupTime) {
				t.Errorf("CreateDate was not assigned")
			}
			if ret.UserID != tc.usrID {
				t.Errorf("User ID mismatch, expect %s, got %s",
					tc.usrID, ret.UserID)
			}
			if !bytes.Equal(ret.Val, tc.key) {
				t.Errorf("API key mismatch, expect %s, got %s",
					tc.key, ret.Val)
			}
			return
		})
	}
}

func TestRoach_APIKeyByUserIDVal(t *testing.T) {
	conf, tearDown := setup(t)
	defer tearDown()
	r := newRoach(t, conf)
	usrID := "123"
	expKey := insertAPIKey(t, r, usrID)
	tt := []struct {
		name        string
		userID      string
		key         []byte
		expNotFound bool
	}{
		{name: "found", userID: usrID, key: expKey.Value(), expNotFound: false},
		{name: "not found key", userID: usrID, key: []byte{}, expNotFound: true},
		{name: "not found userID", userID: "345", key: expKey.Value(), expNotFound: true},
		{name: "not found all", userID: "345", key: []byte{}, expNotFound: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actKey, err := r.APIKeyByUserIDVal(tc.userID, tc.key)
			if tc.expNotFound {
				if !r.IsNotFoundError(err) {
					t.Fatalf("Expected not found error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if !reflect.DeepEqual(expKey, actKey) {
				t.Errorf("API Key mismatch:\nExpect:\t%+v\nGot:\t%+v",
					expKey, actKey)
			}
		})
	}
}

func insertAPIKey(t *testing.T, r *roach.Roach, usrID string) apiH.Key {
	k, err := r.InsertAPIKey(usrID, bytes.Repeat([]byte("x"), 56))
	if err != nil {
		t.Fatalf("Error setting up: insert API key: %v", err)
	}
	return k
}
