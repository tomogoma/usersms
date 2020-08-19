package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/tomogoma/usersms/pkg/bootstrap"
	"github.com/tomogoma/usersms/pkg/config"
	httpIntl "github.com/tomogoma/usersms/pkg/handler/http"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/logging/logrus"
	_ "github.com/tomogoma/usersms/pkg/logging/standard"
)

func main() {

	confFile := flag.String("conf", config.DefaultConfPath(), "location of config file")
	flag.Parse()
	log := &logrus.Wrapper{}
	deps := bootstrap.Instantiate(*confFile, log)

	listenNSrvLg := log.WithField(logging.FieldAction, "Listen and serve")

	port := fmt.Sprintf(":%d", *deps.Config.Service.Port)

	listenNSrvLg.Infof("Will listen on :'%s'", port)

	httpHandler, err := httpIntl.NewHandler(httpIntl.Config{
		Guard:          deps.Guard,
		Logger:         log,
		BaseURL:        config.WebRootPath(),
		Rater:          deps.RatingMan,
		UserProfiler:   deps.UserMan,
		AllowedOrigins: deps.Config.Service.AllowedOrigins,
	})
	logging.LogFatalOnError(listenNSrvLg, err, "Instantiate http Handler")

	logging.LogFatalOnError(
		listenNSrvLg,
		http.ListenAndServe(port, httpHandler),
		"Run server",
	)
}
