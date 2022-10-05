package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "username"
)

func (h *Handler) DeserializeUser(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "invalid auth header"})
		return
	}

	if len(headerParts[1]) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "token is empty"})
		return
	}

	username, err := h.jwt.ValidateToken(headerParts[1], h.cfg.AppConfig.JWTToken.JwtAccessKey)
	if err != nil {
		h.log.Errorln(err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.Set(userCtx, username)

}
