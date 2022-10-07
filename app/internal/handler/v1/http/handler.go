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
		auth := api.Group("auth/")
		{
			auth.POST("register", h.RegisterHandler)
			auth.POST("login", h.Login)
			auth.GET("refresh", h.RefreshAccessToken)
			auth.GET("logout", h.Logout)
			auth.GET("test", h.DeserializeUser, h.Test)
		}
	}

	return r
}

func (h *Handler) Test(ctx *gin.Context) {
	// isAuth, err := h.FetchAuth(*ctx)
	// if !isAuth || err != nil {
	// 	builErrorResponse(ctx, http.StatusBadRequest, Response{
	// 		Status:  statusError,
	// 		Message: "unauthorized",
	// 		Data:    "not cookies username",
	// 	})
	// }
	ctx.JSON(200, "qwertty")
}
