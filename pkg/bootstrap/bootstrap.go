package bootstrap

import (
	"io/ioutil"

	"github.com/sony/sonyflake"
	"github.com/tomogoma/crdb"
	"github.com/tomogoma/go-api-guard"
	"github.com/tomogoma/usersms/pkg/config"
	"github.com/tomogoma/usersms/pkg/db/roach"
	"github.com/tomogoma/usersms/pkg/jwt"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/phone"
	"github.com/tomogoma/usersms/pkg/rating"
	"github.com/tomogoma/usersms/pkg/uid"
	"github.com/tomogoma/usersms/pkg/user"
	"time"
)

type Deps struct {
	Config    config.General
	Guard     *api.Guard
	Roach     *roach.Roach
	JWTEr     *jwt.Manager
	UserMan   *user.Manager
	RatingMan *rating.Manager
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

func InstantiateJWTHandler(lg logging.Logger, srvcConf config.Service) *jwt.Manager {
	JWTKey := []byte(srvcConf.TokenKey)
	if len(JWTKey) == 0 {
		var err error
		JWTKey, err = ioutil.ReadFile(srvcConf.AuthTokenKeyFile)
		logging.LogFatalOnError(lg, err, "Read JWT key file")
	}
	jwter, err := jwt.NewManager(jwt.WithHS256Key(JWTKey))
	logging.LogFatalOnError(lg, err, "Instantiate JWT handler")
	return jwter
}

func Instantiate(confFile string, lg logging.Logger) Deps {

	conf := readConfig(confFile, lg)

	rdb := InstantiateRoach(lg, conf.Database)
	tg := InstantiateJWTHandler(lg, conf.Service)

	g, err := api.NewGuard(rdb, api.WithMasterKey(conf.Service.MasterAPIKey))
	logging.LogFatalOnError(lg, err, "Instantate API access guard")

	idGen := uid.NewSonyFlake(sonyflake.Settings{})

	rater, err := rating.NewManager(tg, rdb, idGen)
	logging.LogFatalOnError(lg, err, "New rating manager")
	go func() {
		for {
			err := rater.SyncUserRatings(conf.Ratings.SyncInterval)
			logging.LogWarnOnError(lg, err, "Sync User Ratings Periodically")
			time.Sleep(conf.Ratings.SyncInterval)
		}
	}()

	userMan, err := user.NewManager(rdb, tg, phone.Formatter{})
	logging.LogFatalOnError(lg, err, "New user manager")

	return Deps{Config: *conf, Guard: g, Roach: rdb, JWTEr: tg, RatingMan: rater, UserMan: userMan}
}

func readConfig(confFile string, lg logging.Logger) *config.General {

	conf := &config.General{}

	if len(confFile) > 0 {
		lg.WithField(logging.FieldAction, "Read config file").Info("started")
		err := config.ReadFile(confFile, conf)
		logging.LogWarnOnError(lg, err, "Read config file")
		lg.WithField(logging.FieldAction, "Read config file").Info("complete")
	}

	lg.WithField(logging.FieldAction, "Read environment config values").Info("started")
	err := config.ReadEnv(conf)
	logging.LogWarnOnError(lg, err, "Read environment config values")
	lg.WithField(logging.FieldAction, "Read environment config values").Info("complete")

	if conf.Service.Port == nil {
		port := 8080
		lg.WithField(logging.FieldAction, "Set default Port").Infof("No port config found fallback to %d", port)
		conf.Service.Port = &port
	}

	return conf
}
