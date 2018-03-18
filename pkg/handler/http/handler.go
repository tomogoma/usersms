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
	"github.com/tomogoma/usersms/pkg/db/queries"
	"github.com/tomogoma/usersms/pkg/logging"
	"github.com/tomogoma/usersms/pkg/rating"
	"github.com/tomogoma/usersms/pkg/user"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type contextKey string

type Guard interface {
	APIKeyValid(key []byte) (string, error)
}

type Rater interface {
	errors.ToHTTPResponser
	RateUser(token, forUserID, comment string, rating int32) error
	Ratings(token string, filter rating.Filter) ([]rating.Rating, error)
}

type UserProfiler interface {
	errors.ToHTTPResponser
	Update(token string, update user.UserUpdate) (*user.User, error)
	User(token, ID string, offsetUpdateDate time.Time) (*user.User, error)
}

type handler struct {
	errors.ErrToHTTP

	guard  Guard
	logger logging.Logger
	rater  Rater
	usrs   UserProfiler
}

const (
	keyAPIKey           = "x-api-key"
	keyUserID           = "userID"
	keyAuthorization    = "Authorization"
	keyOffsetUpdateDate = "offsetUpdateDate"
	keyOffset           = "offset"
	keyCount            = "count"
	keyByUserID         = "byUserID"
	keyForUserID        = "forUserID"
	keyForSection       = "forSection"

	valBearerAuthPrefix = "bearer "

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
	s.handleStatus(r)
	s.handleUserUpdate(r)
	s.handleGetUser(r)
	s.handleRateUser(r)
	s.handleGetRatings(r)
	s.handleDocs(r)
	s.handleNotFound(r)
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
func (s *handler) handleStatus(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/status").
		HandlerFunc(
			s.guardChain(func(w http.ResponseWriter, r *http.Request) {
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
			}),
		)
}

/**
 * @api {PUT} /users/{userID} UpdateUser
 * @apiName Update User Profile
 * @apiVersion 0.1.0
 * @apiGroup Service
 * @apiDescription All declared JSON values are used, including empty strings, except null values.
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (JSON Request Body) {String} [name] New Name.
 * @apiParam (JSON Request Body) {String} [ICEPhone] New (In Case of Emergency) phone number.
 * @apiParam (JSON Request Body) {String="MALE","FEMALE","OTHER"} [gender] New gender.
 * @apiParam (JSON Request Body) {Object} [avatarURL] New profile picture URL.
 * @apiParam (JSON Request Body) {Object} [bio] New brief description of user.
 *
 * @apiUse User200
 *
 */
func (s *handler) handleUserUpdate(r *mux.Router) {
	r.Methods(http.MethodPut).
		PathPrefix("/users/{" + keyUserID + "}").
		HandlerFunc(
			s.guardChain(func(w http.ResponseWriter, r *http.Request) {
				req := struct {
					UserID    string     `json:"userID"`
					Token     string     `json:"token"`
					Name      JSONString `json:"name"`
					ICEPhone  JSONString `json:"ICEPhone"`
					Gender    JSONString `json:"gender"`
					AvatarURL JSONString `json:"avatarURL"`
					Bio       JSONString `json:"bio"`
				}{}

				if err := unmarshalJSONBody(r, &req); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				req.UserID = mux.Vars(r)[keyUserID]

				var err error
				if req.Token, err = getToken(r); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				usr, err := s.usrs.Update(req.Token, user.UserUpdate{
					UserID:    req.UserID,
					Name:      req.Name.ToStringUpdate(),
					ICEPhone:  req.ICEPhone.ToStringUpdate(),
					Gender:    req.Gender.ToStringUpdate(),
					AvatarURL: req.AvatarURL.ToStringUpdate(),
					Bio:       req.Bio.ToStringUpdate(),
				})
				s.respondJsonOn(w, r, req, NewUser(usr), http.StatusOK, err, s.usrs)
			}),
		)
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
 * @apiUse User200
 *
 */
func (s *handler) handleGetUser(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/users/{" + keyUserID + "}").
		HandlerFunc(
			s.guardChain(func(w http.ResponseWriter, r *http.Request) {
				req := struct {
					UserID           string `json:"userID"`
					Token            string `json:"token"`
					OffsetUpdateDate string `json:"offsetUpdateDate"`
				}{
					UserID:           mux.Vars(r)[keyUserID],
					OffsetUpdateDate: r.URL.Query().Get(keyOffsetUpdateDate),
				}

				var err error
				if req.Token, err = getToken(r); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				oud := time.Time{}
				if req.OffsetUpdateDate != "" {
					if oud, err = time.Parse(time.RFC3339, req.OffsetUpdateDate); err != nil {
						err = errors.NewClientf("invalid offsetUpdateDate: %v", err)
						handleError(w, r, req, err, s)
						return
					}
				}

				usr, err := s.usrs.User(req.Token, req.UserID, oud)
				s.respondJsonOn(w, r, req, NewUser(usr), http.StatusOK, err, s.usrs)
			}),
		)
}

/**
 * @api {POST} /ratings/users/{forUserID} RateUser
 * @apiName Rate a user
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiHeader x-api-key the api key
 * @apiHeader Authorization Bearer token containing auth token e.g. "Bearer [value.of.jwt]"
 *
 * @apiParam (URL Param) {String} [forUserID] ID of the user to rate (ratee).
 *
 * @apiParam (JSON Request Body) {Integer{1-5}} rating The rating awarded by rater to ratee.
 * @apiParam (JSON Request Body) {String} [comment] Comment provided by rater.
 *
 * @apiSuccess (200 Response) nil an empty body
 *
 */
func (s *handler) handleRateUser(r *mux.Router) {
	r.Methods(http.MethodPost).
		PathPrefix("/ratings/users/{" + keyForUserID + "}").
		HandlerFunc(
			s.guardChain(func(w http.ResponseWriter, r *http.Request) {
				req := struct {
					Token     string `json:"token"`
					ForUserID string `json:"forUserID"`
					Rating    int32  `json:"rating"`
					Comment   string `json:"comment"`
				}{}

				if err := unmarshalJSONBody(r, &req); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				req.ForUserID = mux.Vars(r)[keyForUserID]

				var err error
				if req.Token, err = getToken(r); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				err = s.rater.RateUser(req.Token, req.ForUserID, req.Token, req.Rating)
				s.respondJsonOn(w, r, req, nil, http.StatusCreated, err, s.rater)
			}),
		)
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
 * @apiParam (URL Query) {String} [byUserID] Filter ratings by rater's userID.
 *		At least one of forUserID or byUserID must be provided.
 * @apiParam (URL Query) {String} [forSection] Filter ratings by section which
 * 		ratee was rated.
 * @apiUse OffsetCount
 *
 * @apiUse RatingsList200
 *
 */
func (s *handler) handleGetRatings(r *mux.Router) {
	r.Methods(http.MethodGet).
		PathPrefix("/ratings/users/{" + keyForUserID + "}").
		HandlerFunc(
			s.guardChain(func(w http.ResponseWriter, r *http.Request) {
				URLQ := r.URL.Query()
				req := struct {
					ForSection string `json:"forSection"`
					ForUserID  string `json:"forUserID"`
					ByUserID   string `json:"byUserID"`
					Token      string `json:"token"`
					Offset     int64  `json:"offset"`
					Count      int32  `json:"count"`
				}{
					ForUserID:  mux.Vars(r)[keyForUserID],
					ByUserID:   URLQ.Get(keyByUserID),
					ForSection: URLQ.Get(keyForSection),
				}

				var err error

				if req.Token, err = getToken(r); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				if req.Offset, err = getOffset(URLQ); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				if req.Count, err = getCount(URLQ); err != nil {
					handleError(w, r, req, err, s)
					return
				}

				rtngs, err := s.rater.Ratings(req.Token, rating.Filter{
					ForUserID:  queries.NewComparisonString(queries.OpET, req.ForUserID),
					ByUserID:   queries.NewComparisonString(queries.OpET, req.ByUserID),
					ForSection: queries.NewComparisonString(queries.OpET, req.ForSection),
					Offset:     req.Offset,
					Count:      req.Count,
				})
				s.respondJsonOn(w, r, req, NewRatings(rtngs), http.StatusOK, err, s.rater)
			}),
		)
}

/**
 * @api {get} /docs Docs
 * @apiName Docs
 * @apiVersion 0.1.0
 * @apiGroup Service
 *
 * @apiSuccess (200) {html} docs Docs page to be viewed on browser.
 *
 */
func (s *handler) handleDocs(r *mux.Router) {
	r.PathPrefix("/" + config.DocsPath).
		Handler(http.FileServer(http.Dir(config.DefaultDocsDir())))
}

func (s handler) handleNotFound(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(
		s.prepLogger(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Nothing to see here", http.StatusNotFound)
		}),
	)
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

func unmarshalJSONBody(r *http.Request, into interface{}) error {

	bodyB, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.NewClientf("unable to read request body: %v", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(bodyB, into); err != nil {
		return errors.NewClientf("invalid request: %v", err)
	}

	return nil
}

func getToken(r *http.Request) (string, error) {

	authHeaders := r.Header[keyAuthorization]
	bearerPrefixLen := len(valBearerAuthPrefix)
	for _, authHeader := range authHeaders {
		if len(authHeader) <= bearerPrefixLen {
			continue
		}

		if strings.HasPrefix(strings.ToLower(authHeader), valBearerAuthPrefix) {
			return authHeader[bearerPrefixLen:], nil
		}
	}
	return "", errors.NewUnauthorizedf("No %stoken was found among the %s headers",
		valBearerAuthPrefix, keyAuthorization)
}

/**
 * @apiDefine OffsetCount
 *
 * @apiParam (URL Query) {Integer} [offset=0] Index from which to fetch(inclusive).
 * @apiParam (URL Query) {Integer} [count=10] Number of items to fetch.
 */

// getOffset extracts offset from r or returns 0 if not found. An error is
// returned if the offset in r is not a valid int64.
func getOffset(r url.Values) (int64, error) {
	offsetStr := r.Get(keyOffset)
	if offsetStr == "" {
		return 0, nil
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		return -1, errors.NewClientf("invalid offset provided: %v", err)
	}
	return offset, nil
}

// getCount extracts count from r or returns 10 if not found. An error is
// returned if the offset in r is not a valid int32.
func getCount(r url.Values) (int32, error) {
	countStr := r.Get(keyCount)
	if countStr == "" {
		return 10, nil
	}
	count, err := strconv.ParseInt(countStr, 10, 32)
	if err != nil {
		return -1, errors.NewClientf("invalid count provided: %v", err)
	}
	return int32(count), nil
}
