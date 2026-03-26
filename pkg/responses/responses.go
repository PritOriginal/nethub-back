package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Success bool       `json:"success"`
	Data    T          `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Message string `json:"message"`
}

func Success[T any](ctx *gin.Context, status int, data T) {
	ctx.JSON(status, Response[T]{
		Success: false,
		Data:    data,
	})
}

func OK[T any](ctx *gin.Context, data T) {
	ctx.JSON(http.StatusOK, Response[T]{
		Success: true,
		Data:    data,
	})
}

func Created[T any](ctx *gin.Context, data T) {
	ctx.JSON(http.StatusCreated, Response[T]{
		Success: true,
		Data:    data,
	})
}

func Fail(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, Response[any]{
		Success: false,
		Error:   &ErrorInfo{Message: message},
	})
}

func Internal(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Response[any]{
		Success: false,
		Error:   &ErrorInfo{Message: message},
	})
}
