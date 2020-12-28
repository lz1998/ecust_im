package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model/user"
	"github.com/lz1998/ecust_im/util"
)

const bearerLength = len("Bearer ")

func CheckLogin(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < bearerLength {
		c.String(http.StatusUnauthorized, "not login")
		c.Abort()
		return
	}
	jwtStr := strings.TrimSpace(authHeader[bearerLength:])
	ecustUser, err := util.JwtParseUser(jwtStr)
	if err != nil {
		c.String(http.StatusUnauthorized, "not login")
		c.Abort()
		return
	}
	c.Set("user", ecustUser)
	c.Next()
}

func Register(c *gin.Context) {
	req := &dto.RegisterReq{}
	if err := c.Bind(req); err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}

	 u, err := user.CreateUser(&user.EcustUser{Password: req.Password, Nickname: req.Nickname})
	if err != nil {
		c.String(http.StatusInternalServerError, "create user error")
		return
	}
	resp := &dto.RegisterResp{
		Ok:       true,
		Msg:      "succeed to register",
		UserInfo: ConvertUserModelToProto(u, true),
	}
	Return(c, resp)
}

func Login(c *gin.Context) {
	req := &dto.LoginReq{}
	if err := c.Bind(req); err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}

	u, err := user.GetUser(req.UserId)
	if err != nil {
		c.String(http.StatusUnauthorized, "user not exists")
		return
	}

	var resp *dto.LoginResp
	if u.Password != req.Password {
		resp = &dto.LoginResp{
			Ok:  false,
			Msg: "password error",
		}
	} else {
		token, err := util.GenerateJwtTokenString(u)
		if err != nil {
			c.String(http.StatusInternalServerError, "generate jwt error")
			return
		}
		resp = &dto.LoginResp{
			Ok:    true,
			Msg:   "login success",
			Token: token,
		}
	}
	Return(c, resp)
}

func SetUserInfo(c *gin.Context) {
	req := &dto.SetUserInfoReq{}
	if err := c.Bind(req); err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	if err := user.UpdateUser(ConvertUsersProtoToModel(req.UserInfos)); err != nil {
		c.String(http.StatusInternalServerError, "update error")
		return
	}
	resp := &dto.SetUserInfoResp{}
	Return(c, resp)
}

func GetUserInfo(c *gin.Context) {
	req := &dto.GetUserInfoReq{}
	if err := c.Bind(req); err != nil {
		c.String(http.StatusBadRequest, "bad request, not protobuf")
		return
	}
	users, err := user.ListUser(req.UserIds)
	if err != nil {
		c.String(http.StatusInternalServerError, "list user error")
		return
	}
	resp := &dto.GetFriendsResp{
		UserInfos: ConvertUsersModelToProto(users, true),
	}
	Return(c, resp)
}

func ConvertUserModelToProto(modelUser *user.EcustUser, ignorePassword bool) *dto.UserInfo {
	return &dto.UserInfo{
		UserId: modelUser.UserId,
		Password: func() string {
			if ignorePassword {
				return ""
			} else {
				return modelUser.Password
			}
		}(),
		Nickname: modelUser.Nickname,
		Email:    "", // TODO 好像没啥用？有空可以做个人信息
	}
}

func ConvertUsersModelToProto(modelUsers []*user.EcustUser, ignorePassword bool) []*dto.UserInfo {
	protoUsers := make([]*dto.UserInfo, 0)
	for _, modelUser := range modelUsers {
		protoUsers = append(protoUsers, ConvertUserModelToProto(modelUser, ignorePassword))
	}
	return protoUsers
}

func ConvertUserProtoToModel(protoUser *dto.UserInfo) *user.EcustUser {
	return &user.EcustUser{
		UserId:   protoUser.UserId,
		Password: protoUser.Password,
		Nickname: protoUser.Nickname,
	}
}

func ConvertUsersProtoToModel(protoUsers []*dto.UserInfo) []*user.EcustUser {
	modelUsers := make([]*user.EcustUser, 0)
	for _, protoUser := range protoUsers {
		modelUsers = append(modelUsers, ConvertUserProtoToModel(protoUser))
	}
	return modelUsers
}
