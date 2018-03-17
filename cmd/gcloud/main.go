package main

import (
	"net/http"

	"github.com/tomogoma/usersms/pkg/bootstrap"
	"github.com/tomogoma/usersms/pkg/config"
	httpInternal "github.com/tomogoma/usersms/pkg/handler/http"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/logging/logrus"
	"google.golang.org/appengine"
)

func main() {

	config.DefaultConfDir("conf")
	log := &logrus.Wrapper{}
	deps := bootstrap.Instantiate(config.DefaultConfPath(), log)

	httpHandler, err := httpInternal.NewHandler(httpInternal.Config{
		Guard:          deps.Guard,
		Logger:         log,
		BaseURL:        config.WebRootPath(),
		Rater:          deps.RatingMan,
		UserProfiler:   deps.UserMan,
		AllowedOrigins: deps.Config.Service.AllowedOrigins,
	})
	logging.LogFatalOnError(log, err, "Instantiate http Handler")

	http.Handle("/", httpHandler)
	appengine.Main()
}
