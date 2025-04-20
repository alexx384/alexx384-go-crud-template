package httputil

import "github.com/gin-gonic/gin"

func NewError(ctx *gin.Context, status int, err error) {
	httpError := HTTPStatusMessage{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, httpError)
}

type HTTPStatusMessage struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
