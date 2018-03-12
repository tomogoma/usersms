package roach_test

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/tomogoma/crdb"
	"github.com/tomogoma/usersms/pkg/db/roach"
	"github.com/tomogoma/usersms/pkg/config"
	"flag"
	"sync/atomic"
)

var (
	confPath = flag.String(
		"conf",
		config.DefaultConfPath(),
		"/path/to/imagems.conf.yml",
	)

	currID = int64(0)
)

func setup(t *testing.T) (crdb.Config, func()) {

	t.Parallel()
	conf, err := config.ReadFile(*confPath)
	if err != nil {
		t.Fatalf("Read config file: %v", err)
	}

	conf.Database.DBName = conf.Database.DBName + "_test_" + strconv.FormatInt(nextID(), 10)

	return conf.Database, func() {
		rdb := getDB(t, conf.Database)
		defer rdb.Close()
		_, err := rdb.Exec("DROP DATABASE " + conf.Database.DBName)
		if err != nil {
			t.Fatalf("Error dropping test db: %v", err)
		}
	}
}

func TestNewRoach(t *testing.T) {

	tt := []struct {
		name   string
		opts   []roach.Option
		expErr bool
	}{
		{
			name: "valid",
			opts: []roach.Option{
				roach.WithDBName("a_test_db_name"),
				roach.WithDSN("a dsn value"),
			},
			expErr: false,
		},
		{
			name:   "valid (no options)",
			opts:   nil,
			expErr: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := roach.NewRoach(tc.opts...)
			if r == nil {
				t.Fatalf("Got nil roach")
			}
		})
	}
}

func TestRoach_InitDBIfNot(t *testing.T) {

	conf, tearDown := setup(t)
	defer tearDown()

	r := newRoach(t, conf)
	rdb := getDB(t, conf)
	defer rdb.Close()
	if err := r.InitDBIfNot(); err != nil {
		t.Fatalf("Initial init call failed: %v", err)
	}

	tt := []struct {
		name       string
		hasVersion bool
		version    []byte
		expErr     bool
	}{
		{
			name:       "first use",
			hasVersion: false,
			expErr:     false,
		},
		{
			name:       "versions equal",
			hasVersion: true,
			version:    []byte(strconv.Itoa(roach.Version)),
			expErr:     false,
		},
		{
			name:       "db version smaller",
			hasVersion: true,
			version:    []byte(strconv.Itoa(roach.Version - 1)),
			expErr:     true,
		},
		{
			name:       "db version bigger",
			hasVersion: true,
			version:    []byte(strconv.Itoa(roach.Version + 1)),
			expErr:     true,
		},
	}

	cols := roach.ColDesc(roach.ColKey, roach.ColValue, roach.ColUpdateDate)
	updCols := roach.ColDesc(roach.ColValue, roach.ColUpdateDate)
	upsertQ := `
		INSERT INTO ` + roach.TblConfigurations + ` (` + cols + `)
			VALUES ('db.version', $1, CURRENT_TIMESTAMP)
			ON CONFLICT (` + roach.ColKey + `)
			DO UPDATE SET (` + updCols + `) = ($1, CURRENT_TIMESTAMP)`
	delQ := `
		DELETE FROM ` + roach.TblConfigurations + `
			WHERE ` + roach.ColKey + `='db.version'`

	for _, tc := range tt {
		if _, err := rdb.Exec(delQ); err != nil {
			t.Fatalf("Error setting up: clear previous config: %v", err)
		}
		if tc.hasVersion {
			if _, err := rdb.Exec(upsertQ, tc.version); err != nil {
				t.Fatalf("Error setting up: insert test config: %v", err)
			}
		}
		t.Run(tc.name, func(t *testing.T) {
			r = newRoach(t, conf)
			err := r.InitDBIfNot()
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				// set db to have correct version (init error should be cached not queried)
				if _, err := rdb.Exec(upsertQ, []byte(strconv.Itoa(roach.Version))); err != nil {
					t.Fatalf("Error setting up: insert test config: %v", err)
				}
				if err := r.InitDBIfNot(); err == nil {
					t.Fatalf("Subsequent init db not returning error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got an error: %v", err)
			}
			// set db to have incorrect version (isInit flag should be cached, not queried)
			if _, err := rdb.Exec(upsertQ, []byte(strconv.Itoa(roach.Version+10))); err != nil {
				t.Fatalf("Error setting up: insert test config: %v", err)
			}
			if err = r.InitDBIfNot(); err != nil {
				t.Fatalf("Subsequent init not working")
			}
		})
	}
}

func newRoach(t *testing.T, conf crdb.Config) *roach.Roach {
	r := roach.NewRoach(
		roach.WithDBName(conf.DBName),
		roach.WithDSN(conf.FormatDSN()),
	)
	if r == nil {
		t.Fatalf("Got nil roach")
	}
	return r
}

func getDB(t *testing.T, conf crdb.Config) *sql.DB {
	DB, err := sql.Open("postgres", conf.FormatDSN())
	if err != nil {
		t.Fatalf("new db instance: %s", err)
	}
	return DB
}

func nextID() int64 {
	return atomic.AddInt64(&currID, 1)
}
