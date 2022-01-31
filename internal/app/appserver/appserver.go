package appserver

import (
	"github.com/gorilla/mux"
	"github.com/ivanmalyi/RestApi/internal/app/store"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type AppServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store *store.Store
}

func New(config *Config)*AppServer {
	return &AppServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (server *AppServer)configureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)
	if err != nil {
		return err
	}
	server.logger.SetLevel(level)

	return nil
}

func (server *AppServer) configureRouter() {
	server.router.HandleFunc("/hello", server.HandleHello())
}

func (server *AppServer) configureStore() error {
	storeConn := store.New(server.config.Store)
	err := storeConn.Open()
	if err != nil {
		return err
	}

	server.store = storeConn

	return nil
}

func (server *AppServer) Start() error {
	err := server.configureLogger()
	if err != nil {
		return err
	}
	server.configureRouter()
	server.logger.Info("server start listen")

	err = server.configureStore()
	if err != nil {
		return err
	}
	return http.ListenAndServe(server.config.BindAddr, server.router)
}

func (server *AppServer) HandleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, "hello world")
	}
}