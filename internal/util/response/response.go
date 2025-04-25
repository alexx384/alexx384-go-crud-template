package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

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

func GetIntParamOrDefault(ctx *gin.Context, paramName string, defaultValue int) (resultVal int, resultErr error) {
	resultVal = defaultValue
	if param, exists := ctx.Params.Get(paramName); exists {
		if value, err := strconv.Atoi(param); err != nil {
			resultErr = err
		} else {
			resultVal = value
		}
	}
	return
}

func GetIntParam(ctx *gin.Context, paramName string) (resultVal int, resultErr error) {
	if param, exists := ctx.Params.Get(paramName); exists {
		if value, err := strconv.Atoi(param); err != nil {
			resultErr = err
		} else {
			resultVal = value
		}
	} else {
		resultErr = fmt.Errorf("param %s not found", paramName)
	}
	return
}
