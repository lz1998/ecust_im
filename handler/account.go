package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model/user"
)

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
		UserInfo: ConvertUserModelToProto(u),
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
		resp = &dto.LoginResp{
			Ok:    true,
			Msg:   "login success",
			Token: "", // TODO 生成jwt token
		}
	}
	Return(c, resp)
}

func SetUserInfo(c *gin.Context) {



}

func GetUserInfo(c *gin.Context) {

}

func ConvertUserModelToProto(modelUser *user.EcustUser) *dto.UserInfo {
	return &dto.UserInfo{
		UserId:   modelUser.UserId,
		Nickname: modelUser.Nickname,
		Email:    "",
	}
}
