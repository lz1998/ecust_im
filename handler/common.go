package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

func Return(c *gin.Context, resp proto.Message) {
	data, err := proto.Marshal(resp)
	if err != nil {
		c.String(http.StatusInternalServerError, "marshal resp error")
		return
	}
	c.Data(http.StatusOK, c.ContentType(), data)
	return
}
