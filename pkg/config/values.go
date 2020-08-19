package config

import (
	"fmt"
	"github.com/Netflix/go-env"
	"io/ioutil"
	"strings"
	"time"

	"github.com/tomogoma/crdb"
	"github.com/tomogoma/go-typed-errors"
	"gopkg.in/yaml.v2"
)

const (
	EnvKeySrvcAllowedOrigins = "SRVC_ALLOWED_ORIGINS"
	EnvKeyDatabaseURL        = "DATABASE_URL"
)

type Service struct {
	RegisterInterval   time.Duration `json:"registerInterval" yaml:"registerInterval" env:"MS_REGISTER_INTERVAL"`
	LoadBalanceVersion string        `json:"loadBalanceVersion" yaml:"loadBalanceVersion" env:"MS_LOAD_BALANCE_VERSION"`
	MasterAPIKey       string        `json:"masterAPIKey" yaml:"masterAPIKey" env:"SRVC_MASTER_API_KEY"`
	AllowedOrigins     []string      `json:"allowedOrigins" yaml:"allowedOrigins" env:"-"`
	AuthTokenKeyFile   string        `json:"authTokenKeyFile" yaml:"authTokenKeyFile"`
	TokenKey           string        `json:"-" yaml:"-" env:"AUTH_JWT_TOKEN_KEY"`
	Port               *int          `json:"port" yaml:"port" env:"PORT"`
}

func unmarshalServcConf(conf *Service) (env.EnvSet, error) {
	es, err := env.UnmarshalFromEnviron(conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if allowedOrigins, exists := es[EnvKeySrvcAllowedOrigins]; exists {
		conf.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}
	return es, nil
}

type Ratings struct {
	SyncInterval time.Duration `json:"syncInterval" yaml:"syncInterval" env:"RATINGS_SYNC_INTERVAL"`
}

type General struct {
	Service     Service     `json:"serviceConfig,omitempty" yaml:"serviceConfig"`
	Database    crdb.Config `json:"database,omitempty" yaml:"database"`
	Ratings     Ratings     `json:"ratings" yaml:"ratings"`
	DatabaseURL string      `json:"databaseURL" yaml:"databaseURL"`
}

func ReadFile(fName string, conf *General) error {
	confD, err := ioutil.ReadFile(fName)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(confD, &conf); err != nil {
		err = errors.Newf("unmarshal conf file (%s) contents: %v",
			fName, err)
		return err
	}
	return nil
}

func ReadEnv(conf *General) error {
	if conf == nil {
		return errors.New("nil config")
	}

	envSet, err := unmarshalServcConf(&conf.Service)
	if err != nil {
		return fmt.Errorf("read service config values: %v", err)
	}
	if err := env.Unmarshal(envSet, &conf.Service); err != nil {
		return fmt.Errorf("read Microservice config values: %v", err)
	}

	if dbURL, exists := envSet[EnvKeyDatabaseURL]; exists {
		conf.DatabaseURL = dbURL
	}

	return nil
}
