package main

import (
	"net"

	"github.com/Asliddin3/customer-servis/config"
	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	"github.com/Asliddin3/customer-servis/kafka"
	"github.com/Asliddin3/customer-servis/pkg/db"
	"github.com/Asliddin3/customer-servis/pkg/logger"
	"github.com/Asliddin3/customer-servis/pkg/messagebroker"
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
		logger.String("database", cfg.PostgresDatabase))
	connDb, err := db.ConnectToDb(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		log.Fatal("error while connect to clients", logger.Error(err))
	}
	publisherMap := make(map[string]messagebroker.Producer)
	customerPublisher := kafka.NewKafkaProducer(cfg, log, "customer.customer")
	defer func() {
		err := customerPublisher.Stop()
		if err != nil {
			log.Fatal("failed to stop kafka costumer", logger.Error(err))
		}
	}()

	publisherMap["customer"] = customerPublisher
	customerService := service.NewCustomerService(grpcClient, connDb, log, publisherMap)
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
