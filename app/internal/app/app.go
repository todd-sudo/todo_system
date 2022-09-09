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
	"github.com/todd-sudo/todo_system/internal/handler"
	"github.com/todd-sudo/todo_system/internal/repository"
	"github.com/todd-sudo/todo_system/internal/service"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"github.com/todd-sudo/todo_system/pkg/server"
)

// logger.Println("swagger docs initializing")
// router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
// router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

// logger.Println("heartbeat metric initializing")
// metricHandler := metric.Handler{}
// metricHandler.Register(router)

func RunApplication() {
	logging.Init()
	log := logging.GetLogger()
	log.Infoln("Connect logger successfully!")

	cfg := config.GetConfig()
	log.Infoln("Connect config successfully!")

	gin.SetMode(cfg.AppConfig.GinMode)

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	db, err := database.NewPostgresDB(cfg, &log)
	if err != nil {
		log.Panicln(err)
	}
	log.Infoln("Connect database successfully!")

	repositories := repository.NewRepository(ctx, db, log)
	log.Info("Connect repositories successfully!")

	services := service.NewService(ctx, *repositories, log)
	log.Info("Connect services successfully!")

	handlers := handler.NewHandler(log, *cfg, services)
	log.Info("Connect services handlers!")

	router := gin.New()

	allFile, err := os.OpenFile("logs/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(fmt.Sprintf("[Message]: %s", err))
	}
	gin.DefaultWriter = io.MultiWriter(allFile)
	router.Use(gin.Logger())

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

// 	c := cors.New(cors.Options{
// 		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
// 		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
// 		AllowCredentials:   true,
// 		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
// 		OptionsPassthrough: true,
// 		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
// 		// Enable Debugging for testing, consider disabling in production
// 		Debug: false,
// 	})

// 	handler := c.Handler(a.router)

// 	a.httpServer = &http.Server{
// 		Handler:      handler,
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}

// 	a.logger.Println("application completely initialized and started")

// 	if err := a.httpServer.Serve(listener); err != nil {
// 		switch {
// 		case errors.Is(err, http.ErrServerClosed):
// 			a.logger.Warn("server shutdown")
// 		default:
// 			a.logger.Fatal(err)
// 		}
// 	}
// 	err := a.httpServer.Shutdown(context.Background())
// 	if err != nil {
// 		a.logger.Fatal(err)
// 	}
// }
