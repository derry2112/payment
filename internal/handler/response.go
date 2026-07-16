package handler

import "github.com/gin-gonic/gin"

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func writeSuccess(c *gin.Context, status int, message string, data any) {
	c.JSON(status, response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func writeError(c *gin.Context, status int, message string, errors any) {
	c.JSON(status, response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
