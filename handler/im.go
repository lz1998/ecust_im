package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model/group"
)

var (
	wsUpgrader = websocket.Upgrader{}
)

func CreateGroup(c *gin.Context) {
	req := &dto.CreateGroupReq{}
	if err := c.Bind(req); err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}

	g, err := group.CreateGroup(&group.EcustGroup{
		GroupName: req.GroupName,
		OwnerId:   0, // TODO jwt做好之后从jwt获取
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "create user error")
		return
	}
	resp := &dto.CreateGroupResp{
		GroupInfo: ConvertGroupModelToProto(g),
	}
	Return(c, resp)
}

func GetFriends(c *gin.Context) {
	// TODO 先做jwt

}

func GetGroups(c *gin.Context) {
	// TODO 先做jwt
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO 先做jwt
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

func ConvertGroupModelToProto(modelGroup *group.EcustGroup) *dto.GroupInfo {
	return &dto.GroupInfo{
		GroupId:   modelGroup.GroupId,
		GroupName: modelGroup.GroupName,
		OwnerId:   modelGroup.OwnerId,
	}
}

func ConvertGroupsModelToProto(modelGroups []*group.EcustGroup) []*dto.GroupInfo {
	protoGroups := make([]*dto.GroupInfo, 0)
	for _, modelGroup := range modelGroups {
		protoGroups = append(protoGroups, ConvertGroupModelToProto(modelGroup))
	}
	return protoGroups
}
