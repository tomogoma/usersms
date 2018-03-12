package main

import (
	"flag"
	"net/http"

	"github.com/micro/go-micro"
	"github.com/micro/go-web"
	"github.com/tomogoma/usersms/pkg/api"
	"github.com/tomogoma/usersms/pkg/bootstrap"
	"github.com/tomogoma/usersms/pkg/config"
	httpIntl "github.com/tomogoma/usersms/pkg/handler/http"
	"github.com/tomogoma/usersms/pkg/handler/rpc"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/logging/logrus"
	_ "github.com/tomogoma/usersms/pkg/logging/standard"
)

func main() {

	confFile := flag.String("conf", config.DefaultConfPath(), "location of config file")
	flag.Parse()
	log := &logrus.Wrapper{}
	deps := bootstrap.Instantiate(*confFile, log)

	serverRPCQuitCh := make(chan error)
	rpcSrv, err := rpc.NewStatusHandler(deps.Guard, log)
	logging.LogFatalOnError(log, err, "Instantate RPC handler")
	go serveRPC(deps.Config.Service, rpcSrv, serverRPCQuitCh)

	serverHttpQuitCh := make(chan error)
	httpHandler, err := httpIntl.NewHandler(deps.Guard, log, config.WebRootPath(),
		deps.Config.Service.AllowedOrigins)
	logging.LogFatalOnError(log, err, "Instantiate HTTP handler")
	go serveHttp(deps.Config.Service, httpHandler, serverHttpQuitCh)

	select {
	case err = <-serverHttpQuitCh:
		logging.LogFatalOnError(log, err, "Serve HTTP")
	case err = <-serverRPCQuitCh:
		logging.LogFatalOnError(log, err, "Serve RPC")
	}
}

func serveRPC(conf config.Service, rpcSrv *rpc.StatusHandler, quitCh chan error) {
	service := micro.NewService(
		micro.Name(config.CanonicalRPCName()),
		micro.Version(conf.LoadBalanceVersion),
		micro.RegisterInterval(conf.RegisterInterval),
	)
	api.RegisterStatusHandler(service.Server(), rpcSrv)
	err := service.Run()
	quitCh <- err
}

func serveHttp(conf config.Service, h http.Handler, quitCh chan error) {
	srvc := web.NewService(
		web.Handler(h),
		web.Name(config.CanonicalWebName()),
		web.Version(conf.LoadBalanceVersion),
		web.RegisterInterval(conf.RegisterInterval),
	)
	quitCh <- srvc.Run()
}
