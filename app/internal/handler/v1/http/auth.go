package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todd-sudo/todo_system/internal/dto"
)

func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var registerDTO dto.InsertUserDTO
	if err := ctx.ShouldBind(&registerDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	h.log.Infof("%+v", registerDTO)
}
