package rpc

import (
	"encoding/json"

	"github.com/pborman/uuid"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/api"
	"github.com/tomogoma/usersms/pkg/config"
	"github.com/tomogoma/usersms/pkg/logging"
	"golang.org/x/net/context"
)

type Guard interface {
	IsUnauthorizedError(error) bool
	IsForbiddenError(error) bool
	APIKeyValid(key []byte) (string, error)
}

type StatusHandler struct {
	errors.NotImplErrCheck
	errors.AuthErrCheck
	errors.ClErrCheck

	guard  Guard
	logger logging.Logger
}

func NewStatusHandler(g Guard, l logging.Logger) (*StatusHandler, error) {
	if g == nil {
		return nil, errors.New("Guard was nil")
	}
	if l == nil {
		return nil, errors.New("Logger was nil")
	}

	return &StatusHandler{guard: g, logger: l}, nil
}

func (sh StatusHandler) prepLogger(method string) logging.Logger {
	log := sh.logger.WithField(logging.FieldTransID, uuid.New())
	log.WithFields(map[string]interface{}{
		logging.FieldRPCMethod:      method,
		logging.FieldRequestHandler: "RPC",
	}).Info("new request")
	return log
}

func (sh *StatusHandler) Check(c context.Context, req *api.Request, resp *api.Response) error {
	log := sh.prepLogger("check")
	_, err := sh.guard.APIKeyValid([]byte(req.APIKey))
	if err != nil {
		reqDataB, _ := json.Marshal(req)
		log = log.WithField(logging.FieldRequest, reqDataB)
		if sh.guard.IsUnauthorizedError(err) {
			log.Warnf("Unauthorized: %v", err)
			return err
		}
		if sh.guard.IsForbiddenError(err) {
			log.Warnf("Forbidden: %v", err)
			return err
		}
		log.Errorf("Error checking API Key Valid (guard): %v", err)
		return errors.Newf("Something wicked happened")
	}
	resp.Name = config.Name
	resp.Version = config.VersionFull
	resp.Description = config.Description
	resp.CanonicalName = config.CanonicalRPCName()
	return nil
}
