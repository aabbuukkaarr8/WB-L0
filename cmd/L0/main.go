package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"

	handlerOrders "L0-arch/internal/api/handler/orders"
	"L0-arch/internal/config"
	"L0-arch/internal/db"
	"L0-arch/internal/kafka"
	repoOrders "L0-arch/internal/repository/orders"
	svcOrders "L0-arch/internal/service/orders"
)

type APIServer struct {
	config *config.Config
	logger *logrus.Logger
	router *gin.Engine
}

func New(config *config.Config) *APIServer {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &APIServer{
		config: config,
		logger: logger,
		router: gin.Default(),
	}

}

func (s *APIServer) Run() error {
	if err := s.configLogger(); err != nil {
		return err
	}

	s.logger.Info("Starting API server")
	return http.ListenAndServe(s.config.Server.Port, s.router)
}

func (s *APIServer) configLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	store := db.New()
	if err := store.Open(cfg.DB.DSN()); err != nil {
		log.Fatal("[db.open]", err)
	}
	defer store.Close()

	k := kafka.NewKafka(cfg.Kafka)
	//repo
	repo := repoOrders.NewRepository(store)

	//service
	svc := svcOrders.NewService(repo)

	//handler
	handler := handlerOrders.NewHandler(svc)

	s := New(cfg)
	s.ConfigureRouter(handler)

	go func() {
		if err := k.StartConsumerGroup(
			cfg.Kafka.Topic,
			cfg.Kafka.GroupID,
			svc,
		); err != nil {
			log.Fatal("kafka consumer:", err)
		}
	}()

	select {}
}
