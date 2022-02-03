package appserver

import (
	"github.com/gorilla/mux"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	logger *logrus.Logger
	router *mux.Router
	store store.Store
}

func NewServer(store store.Store) *server {
	server := &server{
		logger: logrus.New(),
		router: mux.NewRouter(),
		store: store,
	}
	server.configureRouter()

	return server
}

func (server *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.router.ServeHTTP(writer, request)
}

func (server *server) configureRouter() {
	server.router.HandleFunc("/users", server.handleUserCreate()).Methods(http.MethodPost)
}

func (server *server) handleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}