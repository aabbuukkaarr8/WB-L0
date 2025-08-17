package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L0-arch/internal/api/route"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"

	handlerOrders "L0-arch/internal/api/handler/orders"
	"L0-arch/internal/config"
	"L0-arch/internal/db"
	"L0-arch/internal/kafka"
	repoOrders "L0-arch/internal/repository/orders"
	srvOrders "L0-arch/internal/service/orders"
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

	// Создаем основной контекст приложения
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// DB
	store := db.New()
	if err := store.Open(cfg.DB.DSN()); err != nil {
		log.Fatal("[db.open]", err)
	}
	defer store.Close()

	trManager := manager.Must(
		trmsql.NewDefaultFactory(store.GetConn()),
	)

	// Kafka
	k := kafka.New(cfg.Kafka)

	repo := repoOrders.NewRepository(store)
	orderSvc := srvOrders.NewService(repo, trManager)
	h := handlerOrders.NewHandler(orderSvc)

	s := New(cfg)
	route.ConfigureRoutes(s.router, h)
	_ = s.configLogger()

	addr := fmt.Sprintf(cfg.Server.Port)

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	httpErrCh := make(chan error, 1)
	kafkaErrCh := make(chan error, 1)

	go func() {
		s.logger.Infof("HTTP listening on %s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpErrCh <- err
			return
		}
		httpErrCh <- nil
	}()

	go func() {
		kafkaErrCh <- k.StartConsumerGroup(ctx, cfg.Kafka.Topic, cfg.Kafka.GroupID, orderSvc)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down...")

	case err := <-httpErrCh:
		if err != nil {
			s.logger.Errorf("HTTP error: %v", err)
		}

		stop()

	case err := <-kafkaErrCh:
		if err != nil {
			s.logger.Errorf("Kafka error: %v", err)
		}

		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		s.logger.Errorf("HTTP shutdown error: %v", err)
	}

	select {
	case <-kafkaErrCh:
	case <-time.After(3 * time.Second):
	}

	s.logger.Info("Exited cleanly")
}
