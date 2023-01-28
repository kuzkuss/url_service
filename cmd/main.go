package main

import (
	"log"
	"net"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kuzkuss/url_service/cmd/server"
	"github.com/kuzkuss/url_service/config"
	linkDeliveryHttp "github.com/kuzkuss/url_service/internal/link/delivery/http"
	linkDeliveryGrpc "github.com/kuzkuss/url_service/internal/link/delivery/grpc"
	linkRepository "github.com/kuzkuss/url_service/internal/link/repository"
	linkInMem "github.com/kuzkuss/url_service/internal/link/repository/in_memory"
	linkPg "github.com/kuzkuss/url_service/internal/link/repository/postgres"
	linkUsecase "github.com/kuzkuss/url_service/internal/link/usecase"
	link "github.com/kuzkuss/url_service/proto/link"
)

// @title WS Swagger API
// @version 1.0
// @host localhost:8080

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var linkDB linkRepository.RepositoryI

	switch conf.Database {
	case "postgres":
		db, err := gorm.Open(postgres.New(postgres.Config{DSN: conf.PostgresConnectionString}), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		linkDB = linkPg.New(db)
	case "in_memory":
		linkDB = linkInMem.New()
	}

	linkUC := linkUsecase.New(linkDB)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(echoLog.INFO)

	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(echoMiddleware.Recover())

	linkDeliveryHttp.New(e, linkUC)

	lis, err := net.Listen("tcp", conf.HostGRPC + ":" + conf.PortGRPC)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	link .RegisterLinksServer(grpcServer, linkDeliveryGrpc.New(linkUC))

	go func() {
		log.Println("starting server at " + conf.HostGRPC + ":" + conf.PortGRPC)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	s := server.NewServer(e, conf)
	if err := s.Start(conf); err != nil {
		e.Logger.Fatal(err)
	}
}

