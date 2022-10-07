package http

import "github.com/gin-gonic/gin"

const (
	statusOk        = "ok"
	statusError     = "error"
	msgSuccessfully = "successfully"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func builResponse(ctx *gin.Context, statusCode int, response Response) {
	ctx.JSON(statusCode, response)
}

func builErrorResponse(ctx *gin.Context, statusCode int, response Response) {
	ctx.AbortWithStatusJSON(statusCode, response)
}
