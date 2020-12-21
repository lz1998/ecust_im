package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 测试用
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
	return
}
