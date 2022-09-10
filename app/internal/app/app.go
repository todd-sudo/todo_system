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
	apiV1 "github.com/todd-sudo/todo_system/internal/handler/v1/http"
	"github.com/todd-sudo/todo_system/internal/service"
	pgStorage "github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"github.com/todd-sudo/todo_system/pkg/server"
)

func RunApplication() {
	// Init Logger
	logging.Init()
	log := logging.GetLogger()
	log.Infoln("Connect logger successfully!")

	// Init Config
	cfg := config.GetConfig()
	log.Infoln("Connect config successfully!")

	// Init Gin Mode
	gin.SetMode(cfg.AppConfig.GinMode)

	// Init Context
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	// Init Database
	db, err := database.NewPostgresDB(cfg, &log)
	if err != nil {
		log.Panicln(err)
	}
	log.Infoln("Connect database successfully!")

	repositories := pgStorage.NewStorage(ctx, db, log)
	log.Info("Connect repositories successfully!")

	services := service.NewService(ctx, *repositories, log)
	log.Info("Connect services successfully!")

	handlers := apiV1.NewHandler(log, *cfg, services)
	log.Info("Connect services handlers!")

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
			panic("error occurred while running http server: " + err.Error())
		}
	}()
	log.Info("Server started on http://" + cfg.Listen.BindIP + ":" + cfg.Listen.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("Server stopped")

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
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
