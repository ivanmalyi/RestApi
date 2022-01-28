package apiserver

type ApiServer struct {}

func New()*ApiServer {
	return &ApiServer{}
}

func (server *ApiServer) Start() error {
	return nil
}