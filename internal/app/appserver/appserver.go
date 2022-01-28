package appserver

import (
	"github.com/sirupsen/logrus"
)

type AppServer struct {
	config *Config
	logger *logrus.Logger
}

func New(config *Config)*AppServer {


	return &AppServer{
		config: config,
		logger: logrus.New(),
	}
}

func (server *AppServer) Start() error {
	err := server.configureLogger()
	if err != nil {
		return err
	}
	server.logger.Info("server start listen")

	return nil
}

func (server *AppServer)configureLogger() error {
	level, err := logrus.ParseLevel(server.config.LogLevel)
	if err != nil {
		return err
	}
	server.logger.SetLevel(level)

	return nil
}