package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model/friend"
	"github.com/lz1998/ecust_im/model/group"
	"github.com/lz1998/ecust_im/model/group_member"
	"github.com/lz1998/ecust_im/model/request"
	"github.com/lz1998/ecust_im/model/user"
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

func ProcessAdd(c *gin.Context) {
	req := &dto.ProcessAddReq{}
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

	ecustRequest, err := request.GetRequest(req.ReqId)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to get request")
		return
	}

	var processorId int64 // 处理人ID 目标用户/群主
	if ecustRequest.ReqType == request.TFriend {
		processorId = ecustRequest.ToId
	} else {
		ecustGroup, err := group.GetGroup(ecustRequest.ToId)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to get group")
			return
		}
		processorId = ecustGroup.OwnerId
	}

	if processorId != ecustUser.UserId {
		c.String(http.StatusForbidden, "not your request")
		return
	}

	if req.Accept {
		ecustRequest.Status = 1
	} else {
		ecustRequest.Status = 2
	}

	if err := request.UpdateRequest(ecustRequest); err != nil {
		c.String(http.StatusInternalServerError, "failed to update request")
		return
	}

	if req.Accept {
		if ecustRequest.ReqType == 0 { // 好友
			if _, err := friend.SaveFriend(&friend.EcustFriend{
				UserA:  ecustRequest.FromId,
				UserB:  ecustRequest.ToId,
				Status: 1,
			}); err != nil {
				c.String(http.StatusInternalServerError, "failed to save friend")
				return
			}
		} else { // 群
			if _, err := group_member.SaveGroupMember(&group_member.EcustGroupMember{
				GroupId: ecustRequest.ToId,
				UserId:  ecustRequest.FromId,
				Status:  1,
			}); err != nil {
				c.String(http.StatusInternalServerError, "failed to save group member")
				return
			}
		}
	}
	resp := &dto.ProcessAddResp{}
	Return(c, resp)
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
