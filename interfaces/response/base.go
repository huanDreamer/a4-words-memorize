package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Succ(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, BaseResponse{
		Code: http.StatusOK,
		Data: data,
	})
}

func Err(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, BaseResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}
