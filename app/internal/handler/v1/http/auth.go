package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/todd-sudo/todo_system/internal/dto"
)

type registerResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"create_at"`
	Avatar    string    `json:"avatar"`
}

func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var registerDTO dto.InsertUserDTO
	if err := ctx.ShouldBind(&registerDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	user, err := h.service.InsertUser(ctx, &registerDTO)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": "register failed"})
		return
	}
	ctx.JSON(http.StatusOK, &registerResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Avatar:    user.Avatar,
	})
}

func (h *Handler) Login(ctx *gin.Context) {
	var loginDTO dto.InsertUserDTO
	if err := ctx.ShouldBind(&loginDTO); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	user, err := h.service.FindUserByUsername(ctx, loginDTO.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": "you are not registred"})
		return
	}
	if err := h.service.ComparePassword(user.Password, loginDTO.Password); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "error": "wrong password"})
		return
	}

	// Generate Tokens
	accessToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.AccessTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtAccessKey,
	)
	if err != nil {
		h.log.Errorln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.RefreshTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtRefreshKey,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, h.cfg.AppConfig.JWTToken.RefreshTokenMaxAge*60, "/", "localhost", true, true)

	err = h.redisService.SetRefreshToken(ctx, user.Username, refreshToken, time.Duration(h.cfg.AppConfig.JWTToken.RefreshTokenMaxAge))
	if err != nil {
		h.log.Errorln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (h *Handler) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "could not refresh access token"})
		return
	}
	h.log.Infoln(cookie)
	sub, err := h.jwt.ValidateToken(cookie, h.cfg.AppConfig.JWTToken.JwtRefreshKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	user, err := h.service.FindUserByUsername(ctx, sub)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}
	accessToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.AccessTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtAccessKey,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
