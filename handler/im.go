package handler

import (
	"fmt"
	"github.com/lz1998/ecust_im/model/friend"
	"github.com/lz1998/ecust_im/model/group_member"
	"github.com/lz1998/ecust_im/model/user"
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

	tmp, exist := c.Get("user")
	ecustUser := tmp.(*user.EcustUser)
	if !exist {
		c.String(http.StatusUnauthorized, "not login")
		return
	}

	ecustGroup, err := group.CreateGroup(&group.EcustGroup{
		GroupName: req.GroupName,
		OwnerId:   ecustUser.UserId,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "create user error")
		return
	}
	resp := &dto.CreateGroupResp{
		GroupInfo: ConvertGroupModelToProto(ecustGroup),
	}
	Return(c, resp)
}

func GetFriends(c *gin.Context) {
	tmp, exist := c.Get("user")
	ecustUser := tmp.(*user.EcustUser)
	if !exist {
		c.String(http.StatusUnauthorized, "not login")
		return
	}
	friendIds, err := friend.ListFriend(ecustUser.UserId)
	if err != nil {
		c.String(http.StatusInternalServerError, "list friend error")
		return
	}
	users, err := user.ListUser(friendIds)
	if err != nil {
		c.String(http.StatusInternalServerError, "list user error")
		return
	}
	resp := &dto.GetFriendsResp{
		UserInfos: ConvertUsersModelToProto(users, true),
	}
	Return(c, resp)
}

func GetGroups(c *gin.Context) {
	tmp, exist := c.Get("user")
	ecustUser := tmp.(*user.EcustUser)
	if !exist {
		c.String(http.StatusUnauthorized, "not login")
		return
	}
	groupIds, err := group_member.ListGroup(ecustUser.UserId)
	if err != nil {
		c.String(http.StatusInternalServerError, "list group error")
		return
	}
	groups, err := group.ListGroup(groupIds)
	if err != nil {
		c.String(http.StatusInternalServerError, "list group info error")
		return
	}
	resp := &dto.GetGroupsResp{
		GroupInfos: ConvertGroupsModelToProto(groups),
	}
	Return(c, resp)
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
