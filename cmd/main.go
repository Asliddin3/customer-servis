package main

import (
	"net"

	"github.com/Asliddin3/customer-servis/config"
	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	"github.com/Asliddin3/customer-servis/pkg/db"
	"github.com/Asliddin3/customer-servis/pkg/logger"
	"github.com/Asliddin3/customer-servis/service"
	grpcclient "github.com/Asliddin3/customer-servis/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "")
	defer logger.Cleanup(log)

	log.Info("main:sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("datbase", cfg.PostgresDatabase))
	connDb, err := db.ConnectToDb(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		log.Fatal("error while connect to clients", logger.Error(err))
	}

	customerService := service.NewCustomerService(grpcClient, connDb, log)
	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterCustomerServiceServer(s, customerService)
	log.Info("main: server runing",
		logger.String("port", cfg.RPCPort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
