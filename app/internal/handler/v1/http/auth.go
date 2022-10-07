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

const (
	refreshTokenCookie = "refresh_token"
	usernameCookies    = "username"
)

func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var registerDTO dto.InsertUserDTO
	if err := ctx.ShouldBind(&registerDTO); err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "register dto model error",
			Data:    err,
		})
		return
	}
	user, err := h.service.InsertUser(ctx, &registerDTO)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "register failed error",
			Data:    nil,
		})
		return
	}

	builResponse(ctx, http.StatusCreated, Response{
		Status:  statusOk,
		Message: msgSuccessfully,
		Data: &registerResponse{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			Avatar:    user.Avatar,
		},
	})
}

func (h *Handler) Login(ctx *gin.Context) {
	var loginDTO dto.InsertUserDTO
	if err := ctx.ShouldBind(&loginDTO); err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "login dto model error",
			Data:    err,
		})
		return
	}

	user, err := h.service.FindUserByUsername(ctx, loginDTO.Username)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "you are not registred",
			Data:    nil,
		})
		return
	}

	if err := h.service.ComparePassword(user.Password, loginDTO.Password); err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "wrong password",
			Data:    nil,
		})
		return
	}

	// Generate Tokens
	accessToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.AccessTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtAccessKey,
	)
	if err != nil {
		builErrorResponse(ctx, http.StatusConflict, Response{
			Status:  statusError,
			Message: "create access_token error",
			Data:    nil,
		})
		return
	}

	refreshToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.RefreshTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtRefreshKey,
	)
	if err != nil {
		builErrorResponse(ctx, http.StatusConflict, Response{
			Status:  statusError,
			Message: "create refresh_token error",
			Data:    nil,
		})
		return
	}

	// maxAge in seconds * 60 = minutes (60sec * 60sec = 3600sec = 60 minutes)
	ctx.SetCookie(
		refreshTokenCookie,
		refreshToken,
		h.cfg.AppConfig.JWTToken.RefreshTokenMaxAge*60,
		"/",
		h.cfg.AppConfig.Domain,
		true,
		true,
	)
	ctx.Set(userCtx, user.Username)

	if err := h.redisService.SetRefreshToken(
		ctx,
		user.Username,
		refreshToken,
		time.Duration(h.cfg.AppConfig.JWTToken.RefreshTokenMaxAge),
	); err != nil {
		h.log.Errorln(err)
		builErrorResponse(ctx, http.StatusConflict, Response{
			Status:  statusError,
			Message: "set refresh token to redis db error",
			Data:    err,
		})
		return
	}

	builResponse(ctx, http.StatusOK, Response{
		Status:  statusOk,
		Message: msgSuccessfully,
		Data:    map[string]string{"access_token": accessToken},
	})
}

func (h *Handler) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie(refreshTokenCookie)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "could not refresh access token",
			Data:    err,
		})
		return
	}

	sub, err := h.jwt.ValidateToken(cookie, h.cfg.AppConfig.JWTToken.JwtRefreshKey)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "validate token error",
			Data:    nil,
		})
		return
	}

	user, err := h.service.FindUserByUsername(ctx, sub)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "find user by username error",
			Data:    nil,
		})
		return
	}

	accessToken, err := h.jwt.CreateToken(
		time.Duration(h.cfg.AppConfig.JWTToken.AccessTokenExpiresIn),
		user.Username,
		h.cfg.AppConfig.JWTToken.JwtAccessKey,
	)
	if err != nil {
		builErrorResponse(ctx, http.StatusConflict, Response{
			Status:  statusError,
			Message: "create access_token error",
			Data:    nil,
		})
		return
	}

	builResponse(ctx, http.StatusOK, Response{
		Status:  statusOk,
		Message: msgSuccessfully,
		Data:    map[string]string{"access_token": accessToken},
	})
}

func (h *Handler) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(refreshTokenCookie)
	if err != nil {
		builErrorResponse(ctx, http.StatusBadRequest, Response{
			Status:  statusError,
			Message: "unauthorized",
			Data:    "not cookies refresh_token",
		})
		return
	}

	username, err := h.jwt.ValidateToken(refreshToken, h.cfg.AppConfig.JWTToken.JwtRefreshKey)
	if err != nil {
		builErrorResponse(ctx, http.StatusUnauthorized, Response{
			Status:  statusError,
			Message: "unauthorize",
			Data:    nil,
		})
		return
	}

	deleted, err := h.redisService.DelRefreshToken(ctx, username)
	if err != nil || deleted == 0 {
		builErrorResponse(ctx, http.StatusUnauthorized, Response{
			Status:  statusError,
			Message: "unauthorized",
			Data:    nil,
		})
		return
	}

	ctx.SetCookie(refreshTokenCookie, "", -1, "/", h.cfg.AppConfig.Domain, false, true)

	builResponse(ctx, http.StatusOK, Response{
		Status:  statusOk,
		Message: msgSuccessfully,
		Data:    nil,
	})
}
