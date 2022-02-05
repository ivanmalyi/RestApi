package appserver

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	errBadCredentials = errors.New("incorrect email or password")
	sessionName = "ses"
)

type server struct {
	logger *logrus.Logger
	router *mux.Router
	store store.Store
	sessionsStore sessions.Store
}

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
	server.router.HandleFunc("/users", server.handleUserCreate()).Methods(http.MethodPost)
	server.router.HandleFunc("/session", server.handleSessionCreate()).Methods(http.MethodPost)
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