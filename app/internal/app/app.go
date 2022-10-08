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

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"

	jwtToken "github.com/todd-sudo/todo_system/internal/auth/jwt"
	"github.com/todd-sudo/todo_system/internal/config"
	database "github.com/todd-sudo/todo_system/internal/db/postgres"
	"github.com/todd-sudo/todo_system/internal/db/redis"
	apiV1 "github.com/todd-sudo/todo_system/internal/handler/v1/http"
	"github.com/todd-sudo/todo_system/internal/hasher"
	servicePg "github.com/todd-sudo/todo_system/internal/service/postgres"
	serviceRedis "github.com/todd-sudo/todo_system/internal/service/redis"
	pgStorage "github.com/todd-sudo/todo_system/internal/storage/postgres"
	redisStorage "github.com/todd-sudo/todo_system/internal/storage/redis"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"github.com/todd-sudo/todo_system/pkg/server"
)

func RunApplication(saveToFile bool) {

	// Init Logger
	logging.Init(saveToFile)
	log := logging.GetLogger()
	log.Infoln("Connect logger successfully!")

	// Init Config
	cfg := config.GetConfig()
	log.Infoln("Connect config successfully!")

	// Init Context
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	// connect to redis
	redisClient := redis.NewRedisClient(
		ctx,
		&redis.CredentialRedis{
			Host:   cfg.Redis.Host,
			Port:   cfg.Redis.Port,
			Secret: cfg.Redis.Secret,
			Size:   cfg.Redis.Size,
		},
		log,
	)
	rc, err := redisClient.ConnectToRedis()
	if err != nil {
		log.Panicln("error connecting to redis %w", err)
	}

	// Init redis store and cookies
	// redisStore, err := redisClient.GetStore()
	// if err != nil {
	// 	log.Panicln("error connecting to redis store %w", err)
	// }
	// redisStore.Options(sessions.Options{
	// 	Secure:   true,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteStrictMode,
	// })
	// log.Infoln("Connect redis successfully!")

	// Init Database
	db, err := database.NewPostgresDB(cfg, &log)
	if err != nil {
		log.Panicln(err)
	}
	log.Infoln("Connect database successfully!")

	// Connect JWT Token
	jwt := jwtToken.NewJWTToken(log, *cfg)
	log.Infoln("Connect JWT Token successfully!")

	// Init hasher password
	hasher := hasher.NewSHA1Hasher(log)
	log.Infoln("Connect hasher successfully!")

	// Connect storages
	storagePg := pgStorage.NewStorage(ctx, db, log)
	log.Infoln("Connect storage postgres successfully!")

	storageRedis := redisStorage.NewJWTStorage(ctx, rc)
	log.Infoln("Connect storage redis successfully!")

	// Connect services
	servicesPG := servicePg.NewService(ctx, *storagePg, log, hasher)
	log.Infoln("Connect service postgres successfully!")

	servicesRedis := serviceRedis.NewRedisService(ctx, rc, storageRedis)
	log.Infoln("Connect service redis successfully!")

	// Connect handlers
	handlers := apiV1.NewHandler(log, *cfg, servicesPG, jwt, servicesRedis)
	log.Infoln("Connect services handlers!")

	// New Gin router
	router := gin.New()
	// Init Gin Mode
	gin.SetMode(cfg.AppConfig.GinMode)
	// router.Use(sessions.Sessions(cfg.AppConfig.Auth.SessionName, redisStore))
	// log.Infoln("Connect redis to GIN successfully")

	// Gin Logs
	enableGinLogs(saveToFile, router)

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
		AllowedOrigins:     []string{"http://127.0.0.1:8000", "http://127.0.0.1:8000", "http://localhost:8000"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(handlers.InitRoutes(router))
	return handler
}

// enableGinLogs включает/отключает gin логи
func enableGinLogs(saveToFile bool, router *gin.Engine) {
	if saveToFile {
		allFile, err := os.OpenFile("logs/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			panic(fmt.Sprintf("[Message]: %s", err))
		}
		gin.DefaultWriter = io.MultiWriter(allFile)
	}

	router.Use(gin.Logger())
}
