package config

import (
	"io/ioutil"
	"time"

	"github.com/tomogoma/crdb"
	"github.com/tomogoma/go-typed-errors"
	"gopkg.in/yaml.v2"
)

type Service struct {
	RegisterInterval   time.Duration `json:"registerInterval,omitempty" yaml:"registerInterval"`
	LoadBalanceVersion string        `json:"loadBalanceVersion,omitempty" yaml:"loadBalanceVersion"`
	MasterAPIKey       string        `json:"masterAPIKey,omitempty" yaml:"masterAPIKey"`
	AllowedOrigins     []string      `json:"allowedOrigins" yaml:"allowedOrigins"`
	AuthTokenKeyFile   string        `json:"authTokenKeyFile" yaml:"authTokenKeyFile"`
}

type Auth struct {
}

type General struct {
	Service  Service     `json:"serviceConfig,omitempty" yaml:"serviceConfig"`
	Database crdb.Config `json:"database,omitempty" yaml:"database"`
}

func ReadFile(fName string) (conf General, err error) {
	confD, err := ioutil.ReadFile(fName)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(confD, &conf); err != nil {
		err = errors.Newf("unmarshal conf file (%s) contents: %v",
			fName, err)
		return
	}
	return
}
