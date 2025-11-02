package controller

import (
	"fmt"
	"net/http"

	"github.com/Oleska1601/WBDelayedNotifier/config"
)

type Server struct {
	Srv     *http.Server
	usecase UsecaseInterface
}

func New(cfg *config.ServerConfig, usecase UsecaseInterface) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &Server{
		Srv: &http.Server{
			Addr: addr,
		},
		usecase: usecase,
	}

	server.setupRouter()
	return server
}
