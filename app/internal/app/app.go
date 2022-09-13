package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"github.com/todd-sudo/todo_system/internal/config"
	database "github.com/todd-sudo/todo_system/internal/db/postgres"
	"github.com/todd-sudo/todo_system/internal/db/redis"
	apiV1 "github.com/todd-sudo/todo_system/internal/handler/v1/http"
	"github.com/todd-sudo/todo_system/internal/hasher"
	service_pg "github.com/todd-sudo/todo_system/internal/service/postgres"
	pgStorage "github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"github.com/todd-sudo/todo_system/pkg/server"
)

func RunApplication() {
	// Init Context
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	// Init Logger
	logging.Init()
	log := logging.GetLogger()
	log.Infoln("Connect logger successfully!")

	// Init Config
	cfg := config.GetConfig()
	log.Infoln("Connect config successfully!")

	// connect to redis
	rc, err := redis.NewRedisClient(
		ctx,
		&redis.CredentialRedis{Host: cfg.Redis.Host, Port: cfg.Redis.Port},
		log,
	).ConnectToRedis()
	if err != nil {
		log.Panicln("error connecting to redis %w", err)
	}
	log.Infoln("Connect redis successfully!")

	// Init Gin Mode
	gin.SetMode(cfg.AppConfig.GinMode)

	// Init Database
	db, err := database.NewPostgresDB(cfg, &log)
	if err != nil {
		log.Panicln(err)
	}
	log.Infoln("Connect database successfully!")

	// Init hasher password
	hasher := hasher.NewSHA1Hasher(cfg.AppConfig.Auth.PasswordHashSalt)
	log.Infoln("Connect hasher successfully!")

	storages := pgStorage.NewStorage(ctx, db, log)
	log.Infoln("Connect repositories successfully!")

	servicesPG := service_pg.NewService(ctx, *storages, log, hasher)
	log.Infoln("Connect services successfully!")

	handlers := apiV1.NewHandler(log, *cfg, servicesPG)
	log.Infoln("Connect services handlers!")

	// New Gin router
	router := gin.New()

	// Gin Logs
	enableGinLogs(true, router)

	// Init Routes and CORS
	handler := initRoutesAndCORS(router, handlers)

	// Start HTTP Server
	srv := server.NewServer(cfg.Listen.Port, handler)
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Panicln("error occurred while running http server: " + err.Error())
		}
	}()
	log.Infoln("Server started on http://" + cfg.Listen.BindIP + ":" + cfg.Listen.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("Server stopped")

	if err := srv.Stop(ctx); err != nil {
		log.Panicf("failed to stop server: %v\n", err)
	}

	if err := rc.Close(); err != nil {
		log.Panicf("error closing Redis Client: %w\n", err)
	}
}

// initRoutesAndCORS инициализирует роутер и обработчики
func initRoutesAndCORS(router *gin.Engine, handlers *apiV1.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:8000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(handlers.InitRoutes(router))
	return handler
}

// enableGinLogs включает/отключает gin логи
func enableGinLogs(enable bool, router *gin.Engine) {
	if enable {
		allFile, err := os.OpenFile("logs/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			panic(fmt.Sprintf("[Message]: %s", err))
		}
		gin.DefaultWriter = io.MultiWriter(allFile)
		router.Use(gin.Logger())
	}
}
