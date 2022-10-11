package http

import (
	"github.com/gin-gonic/gin"
	"github.com/todd-sudo/todo_system/internal/auth/jwt"
	"github.com/todd-sudo/todo_system/internal/config"
	service_pg "github.com/todd-sudo/todo_system/internal/service/postgres"
	redisService "github.com/todd-sudo/todo_system/internal/service/redis"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Handler struct {
	service      *service_pg.Service
	cfg          config.Config
	log          logging.Logger
	jwt          jwt.JWTToken
	redisService redisService.RedisService
}

func NewHandler(
	log logging.Logger,
	cfg config.Config,
	service *service_pg.Service,
	jwt jwt.JWTToken,
	redisService redisService.RedisService,
) *Handler {
	return &Handler{
		service:      service,
		cfg:          cfg,
		log:          log,
		jwt:          jwt,
		redisService: redisService,
	}
}

func (h *Handler) InitRoutes(r *gin.Engine) *gin.Engine {

	api := r.Group("api/")
	{
		api.GET("test", h.DeserializeUser, h.Test)

		auth := api.Group("auth/")
		{
			auth.POST("register", h.RegisterHandler)
			auth.POST("login", h.Login)
			auth.GET("refresh", h.RefreshAccessToken)
			auth.GET("logout", h.Logout)
		}

		folder := api.Group("folder/")
		{
			folder.POST("all")
			folder.POST("create")
			folder.PATCH("update")
			folder.DELETE("delete")
		}

		item := api.Group("item/")
		{
			item.POST("all-folder")
			item.POST("all")
			item.POST("create")
			item.PATCH("update")
			item.DELETE("delete")
		}
	}

	return r
}

func (h *Handler) Test(ctx *gin.Context) {
	ctx.JSON(200, "qwerty")
}
