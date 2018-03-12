package bootstrap

import (
	"io/ioutil"

	"github.com/tomogoma/go-api-guard"
	"github.com/tomogoma/jwt"
	"github.com/tomogoma/usersms/pkg/config"
	"github.com/tomogoma/usersms/pkg/db/roach"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/crdb"
)

type Deps struct {
	Config config.General
	Guard  *api.Guard
	Roach  *roach.Roach
	JWTEr  *jwt.Handler
}

func InstantiateRoach(lg logging.Logger, conf crdb.Config) *roach.Roach {
	var opts []roach.Option
	if dsn := conf.FormatDSN(); dsn != "" {
		opts = append(opts, roach.WithDSN(dsn))
	}
	if dbn := conf.DBName; dbn != "" {
		opts = append(opts, roach.WithDBName(dbn))
	}
	rdb := roach.NewRoach(opts...)
	err := rdb.InitDBIfNot()
	logging.LogWarnOnError(lg, err, "Initiate Cockroach DB connection")
	return rdb
}

func InstantiateJWTHandler(lg logging.Logger, tknKyF string) *jwt.Handler {
	JWTKey, err := ioutil.ReadFile(tknKyF)
	logging.LogFatalOnError(lg, err, "Read JWT key file")
	jwter, err := jwt.NewHandler(JWTKey)
	logging.LogFatalOnError(lg, err, "Instantiate JWT handler")
	return jwter
}

func Instantiate(confFile string, lg logging.Logger) Deps {

	conf, err := config.ReadFile(confFile)
	logging.LogFatalOnError(lg, err, "Read config file")

	rdb := InstantiateRoach(lg, conf.Database)
	tg := InstantiateJWTHandler(lg, conf.Service.AuthTokenKeyFile)

	g, err := api.NewGuard(rdb, api.WithMasterKey(conf.Service.MasterAPIKey))
	logging.LogFatalOnError(lg, err, "Instantate API access guard")

	return Deps{Config: conf, Guard: g, Roach: rdb, JWTEr: tg}
}
