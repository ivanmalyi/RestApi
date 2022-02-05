package appserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	sessionName = "ses"
	CtxKeyUser ctxKey = iota
	CtxKeyRequestID ="XRequestID"
)

var (
	errBadCredentials = errors.New("incorrect email or password")
	errNotAuth = errors.New("not authenticated")
)

type server struct {
	logger *logrus.Logger
	router *mux.Router
	store store.Store
	sessionsStore sessions.Store
}

type ctxKey int8

func NewServer(store store.Store, sessionsStore sessions.Store) *server {
	server := &server{
		logger: logrus.New(),
		router: mux.NewRouter(),
		store: store,
		sessionsStore: sessionsStore,
	}
	server.configureRouter()

	return server
}

func (server *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.router.ServeHTTP(writer, request)
}

func (server *server) configureRouter() {
	server.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	server.router.Use(server.setRequestID)
	server.router.Use(server.logRequest)

	server.router.HandleFunc("/users", server.handleUserCreate()).Methods(http.MethodPost)
	server.router.HandleFunc("/session", server.handleSessionCreate()).Methods(http.MethodPost)

	privateRouter := server.router.PathPrefix("/private").Subrouter()
	privateRouter.Use(server.authenticateUser)
	privateRouter.HandleFunc("/whoami", server.handleWhoami()).Methods(http.MethodGet)
}

func (server *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		logger := server.logger.WithFields(logrus.Fields{
			"remote_addr": req.RemoteAddr,
			"request_id": req.Context().Value(CtxKeyRequestID),
		})

		logger.Infof("started: %s %s", req.Method, req.RequestURI)
		rw := &ResponseWriter{writer, http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rw, req)
		logger.Infof(
			"complited with %d %s %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (server *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-request-ID", id)
		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), CtxKeyRequestID, id)))
	})
}

func (server *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		session, err := server.sessionsStore.Get(req, sessionName)
		if err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			server.error(writer, req, http.StatusUnauthorized, errNotAuth)
			return
		}

		user, err := server.store.User().Find(id.(int))
		if err != nil {
			server.error(writer, req, http.StatusUnauthorized, errNotAuth)
			return
		}

		next.ServeHTTP(writer, req.WithContext(context.WithValue(req.Context(), CtxKeyUser, user)))
	})
}

func (server *server) handleWhoami() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		server.respond(writer, req, http.StatusOK, req.Context().Value(CtxKeyUser).(*model.User))
	}
}

func (server *server) handleUserCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		request := &request{}
		err := json.NewDecoder(req.Body).Decode(request)
		if err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}
		user := &model.User{
			Email: request.Email,
			Password: request.Password,
		}

		err = server.store.User().Create(user)
		if err != nil {
			server.error(writer, req, http.StatusUnprocessableEntity, err)
			return
		}

		user.Sanitize()
		server.respond(writer, req, http.StatusCreated, user)
	}
}

func (server *server) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		request := &request{}
		err := json.NewDecoder(req.Body).Decode(request)
		if err != nil {
			server.error(writer, req, http.StatusBadRequest, err)
			return
		}

		user, err := server.store.User().FindByEmail(request.Email)
		if err != nil || !user.ComparePassword(request.Password) {
			server.error(writer, req, http.StatusUnauthorized, errBadCredentials)
			return
		}

		session, err := server.sessionsStore.Get(req, sessionName)
		if err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = user.ID
		err = server.sessionsStore.Save(req, writer, session)
		if err != nil {
			server.error(writer, req, http.StatusInternalServerError, err)
			return
		}

		server.respond(writer, req, http.StatusOK, nil)
	}
}

func (server *server) error(writer http.ResponseWriter, req *http.Request, code int, err error) {
	server.respond(writer, req, code, map[string]string{"error":err.Error()})
}

func (server *server) respond(writer http.ResponseWriter, req *http.Request, code int, data interface{}) {
	writer.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(writer).Encode(data)
	}
}