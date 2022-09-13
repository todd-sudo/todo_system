package http

import (
	"github.com/gin-gonic/gin"
	"github.com/todd-sudo/todo_system/internal/config"
	service_pg "github.com/todd-sudo/todo_system/internal/service/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
)

type Handler struct {
	service *service_pg.Service
	cfg     config.Config
	log     logging.Logger
}

func NewHandler(log logging.Logger, cfg config.Config, service *service_pg.Service) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
		log:     log,
	}
}

func (h *Handler) InitRoutes(r *gin.Engine) *gin.Engine {
	// var c *gin.Context
	// // c.SetSameSite(http.SameSiteStrictMode)
	// c.SetCookie()
	return r
}
