package server

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
	"github.com/kuzkuss/url_service/config"
)

type Server struct {
	http.Server
}

func NewServer(e *echo.Echo, conf *config.Config) *Server {
	return &Server{
		http.Server{
			Addr:              conf.HostHTTP + ":" + conf.PortHTTP,
			Handler:           e,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
	}
}

func (s *Server) Start(conf *config.Config) error {
	log.Println("starting server at " + conf.HostHTTP + ":" + conf.PortHTTP)
	return s.ListenAndServe()
}

