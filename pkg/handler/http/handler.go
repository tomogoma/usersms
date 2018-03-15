package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/tomogoma/go-typed-errors"
	"github.com/tomogoma/usersms/pkg/config"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/user"
	"github.com/tomogoma/usersms/pkg/rating"
)

type contextKey string

type Guard interface {
	APIKeyValid(key []byte) (string, error)
}

type Rater interface {
	errors.ToHTTPResponser
	RateUser(token string, rating rating.Rating) error
	Ratings(token string, filter rating.Filter) ([]rating.Rating, error)
}

type UserProfiler interface {
	errors.ToHTTPResponser
	Update(token string, update user.UserUpdate) (*user.User, error)
	User(token, ID, offsetUpdateDate string) (*user.User, error)
}

type handler struct {
	errors.ErrToHTTP

	guard  Guard
	logger logging.Logger
	rater  Rater
	usrs   UserProfiler
}

const (
	keyAPIKey = "x-api-key"

	ctxKeyLog = contextKey("log")
)

func NewHandler(conf Config) (http.Handler, error) {

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	r := mux.NewRouter().PathPrefix(conf.BaseURL).Subrouter()
	handler{guard: conf.Guard, logger: conf.Logger, rater: conf.Rater, usrs: conf.UserProfiler}.
		handleRoute(r)

	corsOpts := []handlers.CORSOption{
		handlers.AllowedHeaders([]string{
			"X-Requested-With", "Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization", "X-api-key",
		}),
		handlers.AllowedOrigins(conf.AllowedOrigins),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	}
	return handlers.CORS(corsOpts...)(r), nil
}

func (s handler) handleRoute(r *mux.Router) {

	r.PathPrefix("/status").
		Methods(http.MethodGet).
		HandlerFunc(s.guardChain(s.handleStatus))

	r.PathPrefix("/" + config.DocsPath).
		Handler(http.FileServer(http.Dir(config.DefaultDocsDir())))

	r.NotFoundHandler = http.HandlerFunc(s.prepLogger(s.handleNotFound))
}

/**
 * @api {get} /status Status
 * @apiName Status
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 *
 * @apiSuccess (200) {String} name Micro-service name.
 * @apiSuccess (200)  {String} version http://semver.org version.
 * @apiSuccess (200)  {String} description Short description of the micro-service.
 * @apiSuccess (200)  {String} canonicalName Canonical name of the micro-service.
 *
 */
func (s *handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.respondJsonOn(w, r, nil, struct {
		Name          string `json:"name"`
		Version       string `json:"version"`
		Description   string `json:"description"`
		CanonicalName string `json:"canonicalName"`
	}{
		Name:          config.Name,
		Version:       config.VersionFull,
		Description:   config.Description,
		CanonicalName: config.CanonicalWebName(),
	}, http.StatusOK, nil, s)
}

/**
 * @api {PUT} /users/{userID} UpdateUser
 * @apiName Update User Profile
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (JSON Body) {String} [name] New Name
 * @apiParam (JSON Body) {String} [ICEPhone] New (In Case of Emergency)
 * 		phone number
 * @apiParam (JSON Body) {String="MALE","FEMALE","OTHER"} [gender] New gender
 * @apiParam (JSON Body) {Object} [avatarURL] New profile picture URL
 * @apiParam (JSON Body) {Object} [bio] New brief description of user
 *
 * @apiSuccess (200) {JSON} body  Updated <a href="#api-Objects-User">user</a> object
 *
 */
func (s *handler) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
}

/**
 * @api {GET} /users/{userID} GetUser
 * @apiName Get user
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (URL Param) {String} [userID] ID of the user to fetch.
 *
 * @apiParam (URL Query) {Integer} [offsetUpdateDate] Earliest ISO8601 date that
 * 		the user should have been updated. If the userID exists but the update
 *		date is earlier than this value then a 404 will be returned.
 *
 * @apiSuccess (200) {JSON} body  Updated <a href="#api-Objects-User">user</a> object
 *
 */
func (s *handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
}

/**
 * @api {POST} /ratings/users/{userID} RateUser
 * @apiName Rate a user
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (URL Param) {String} [userID] ID of the user to rate (ratee).
 *
 * @apiParam (JSON Body) {String} byUserID ID of the user rating (rater).
 * @apiParam (JSON Body) {Integer{0-5}} rating The rating awarded by rater to ratee.
 * @apiParam (JSON Body) {String} [comment] Comment provided by rater.
 *
 * @apiSuccess (200) 200
 *
 */
func (s *handler) handleRateUser(w http.ResponseWriter, r *http.Request) {
}

/**
 * @api {GET} /ratings/users/{forUserID} GetRatingsOnUser
 * @apiName Get Ratings On User
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (URL Param) {String} [forUserID] Filter ratings by ratee's userID.
 *		At least one of forUserID or byUserID must be provided.
 *
 * @apiParam (URL Query) {Integer} [offset=0] Index from which to fetch ratings (inclusive).
 * @apiParam (URL Query) {Integer} [count=10] Number of ratings to fetch.
 * @apiParam (URL Query) {String} [byUserID] Filter ratings by rater's userID.
 *		At least one of forUserID or byUserID must be provided.
 *
 * @apiSuccess (200) {JSON} body Array of <a href="#api-Objects-Rating">ratings</a>.
 *
 */
func (s *handler) handleGetRatings(w http.ResponseWriter, r *http.Request) {
}

func (s handler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Nothing to see here", http.StatusNotFound)
}

func (s *handler) guardChain(next http.HandlerFunc) http.HandlerFunc {
	return s.prepLogger(s.guardRoute(next))
}

func (s handler) prepLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := s.logger.WithHTTPRequest(r).
			WithField(logging.FieldTransID, uuid.New())

		log.WithFields(map[string]interface{}{
			logging.FieldURLPath:    r.URL.Path,
			logging.FieldHTTPMethod: r.Method,
		}).Info("new request")

		ctx := context.WithValue(r.Context(), ctxKeyLog, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *handler) guardRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey := r.Header.Get(keyAPIKey)
		clUsrID, err := s.guard.APIKeyValid([]byte(APIKey))
		log := r.Context().Value(ctxKeyLog).(logging.Logger).
			WithField(logging.FieldClientAppUserID, clUsrID)
		ctx := context.WithValue(r.Context(), ctxKeyLog, log)
		if err != nil {
			handleError(w, r.WithContext(ctx), nil, err, s)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// respondJsonOn marshals respData to json and writes it and the code as the
// http header to w. If err is not nil, handleError is called instead of the
// documented write to w.
func (s *handler) respondJsonOn(w http.ResponseWriter, r *http.Request, reqData interface{},
	respData interface{}, code int, err error, errSrc errors.ToHTTPResponser) int {

	if err != nil {
		handleError(w, r, reqData, err, errSrc)
		return 0
	}

	respBytes, err := json.Marshal(respData)
	if err != nil {
		handleError(w, r, reqData, err, errSrc)
		return 0
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	i, err := w.Write(respBytes)
	if err != nil {
		log := r.Context().Value(ctxKeyLog).(logging.Logger)
		log.Errorf("unable write data to response stream: %v", err)
		return i
	}

	return i
}

// handleError writes an error to w using errSrc's logic and logs the error
// using the logger acquired by the prepLogger middleware on r. reqData is
// included in the log data.
func handleError(w http.ResponseWriter, r *http.Request, reqData interface{}, err error, errSrc errors.ToHTTPResponser) {
	reqDataB, _ := json.Marshal(reqData)
	log := r.Context().Value(ctxKeyLog).(logging.Logger).
		WithField(logging.FieldRequest, string(reqDataB))

	if code, ok := errSrc.ToHTTPResponse(err, w); ok {
		log.WithField(logging.FieldResponseCode, code).Warn(err)
		return
	}

	log.WithField(logging.FieldResponseCode, http.StatusInternalServerError).
		Error(err)
	http.Error(w, "Something wicked happened, please try again later",
		http.StatusInternalServerError)
}
