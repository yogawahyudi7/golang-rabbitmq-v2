package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"golang-rabbitmq-v2/internal/config"
	"golang-rabbitmq-v2/internal/database"
	"golang-rabbitmq-v2/internal/domain"
	"golang-rabbitmq-v2/internal/handler"
	"golang-rabbitmq-v2/internal/middleware"
	"golang-rabbitmq-v2/internal/repository"
	"golang-rabbitmq-v2/internal/usecase"
	"golang-rabbitmq-v2/pkg/logger"
	"golang-rabbitmq-v2/pkg/rabbitmq"
	"golang-rabbitmq-v2/pkg/validator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.New(&cfg.Logger)
	log.WithContext("main").Info("Starting application")

	db, err := database.NewPostgresConnection(&cfg.Database, log)
	if err != nil {
		log.WithContext("main").WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	if err := db.Migrate(&domain.User{}, &domain.Message{}); err != nil {
		log.WithContext("main").WithError(err).Fatal("Failed to migrate database")
	}

	rabbitConn, err := rabbitmq.NewConnection(&cfg.RabbitMQ, log)
	if err != nil {
		log.WithContext("main").WithError(err).Fatal("Failed to connect to RabbitMQ")
	}
	defer rabbitConn.Close()

	publisher, err := rabbitmq.NewPublisher(rabbitConn, log)
	if err != nil {
		log.WithContext("main").WithError(err).Fatal("Failed to create publisher")
	}
	defer publisher.Close()

	consumer := rabbitmq.NewConsumer(rabbitConn, log, 10)

	metrics := rabbitConn.GetMetrics()
	metrics.StartPeriodicLogging(30 * time.Second)

	validator := validator.New(log)

	userRepo := repository.NewUserRepository(db, log)
	messageRepo := repository.NewMessageRepository(db, log)

	userUseCase := usecase.NewUserUseCase(userRepo, publisher, validator, log)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, userRepo, publisher, validator, log)

	userHandler := handler.NewUserHandler(userUseCase, log)
	messageHandler := handler.NewMessageHandler(messageUseCase, rabbitConn, log)
	monitoringHandler := handler.NewMonitoringHandler(rabbitConn, log)

	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	router.Use(middleware.Recovery(log))
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.CORS())

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("", userHandler.ListUsers)
		}

		messages := api.Group("/messages")
		{
			messages.POST("", messageHandler.CreateMessage)
			messages.GET("/:id", messageHandler.GetMessage)
			messages.GET("", messageHandler.ListMessages)
			messages.POST("/retry", messageHandler.RetryFailedMessages)
			messages.GET("/stats", messageHandler.GetMessageStats)
		}

		monitoring := api.Group("/monitoring")
		{
			// Basic monitoring
			monitoring.GET("/health", monitoringHandler.GetHealthCheck)
			monitoring.GET("/metrics", monitoringHandler.GetMetrics)
			monitoring.GET("/metrics/detailed", monitoringHandler.GetDetailedMetrics)

			// Real-time activity monitoring
			monitoring.GET("/activity", monitoringHandler.GetActivitySnapshot)

			// Event history
			monitoring.GET("/events", monitoringHandler.GetEventHistory)
			monitoring.GET("/events/publish", monitoringHandler.GetPublishEvents)
			monitoring.GET("/events/consume", monitoringHandler.GetConsumeEvents)
			monitoring.GET("/events/failed", monitoringHandler.GetFailedEvents)
			monitoring.DELETE("/events", monitoringHandler.ClearEventHistory)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		messageHandler := func(ctx context.Context, body []byte) error {
			var queueMessage domain.QueueMessage
			if err := json.Unmarshal(body, &queueMessage); err != nil {
				log.WithContext("consumer").WithError(err).Error("Failed to unmarshal queue message")
				return err
			}

			return messageUseCase.ProcessMessage(ctx, &queueMessage)
		}

		if err := consumer.Start(ctx, messageHandler); err != nil {
			log.WithContext("consumer").WithError(err).Error("Consumer stopped with error")
		}
	}()

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.WithContext("main").WithFields(map[string]interface{}{
			"port": cfg.Server.Port,
			"mode": cfg.Server.GinMode,
		}).Info("Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithContext("main").WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.WithContext("main").Info("Shutting down server...")

	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithContext("main").WithError(err).Error("Server forced to shutdown")
	} else {
		log.WithContext("main").Info("Server exited gracefully")
	}

	if err := consumer.Stop(); err != nil {
		log.WithContext("main").WithError(err).Error("Failed to stop consumer")
	}

	log.WithContext("main").Info("Application stopped")
}