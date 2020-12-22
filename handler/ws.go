package handler

import (
	"net/http"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model/user"
	"github.com/lz1998/ecust_im/util"
	log "github.com/sirupsen/logrus"
)

var (
	wsUpgrader = websocket.Upgrader{}
)

type SendingMessage struct {
	MessageType int
	Data        []byte
}

type UserSession struct {
	Conn        *websocket.Conn
	User        *user.EcustUser
	SendChannel chan *SendingMessage
}

var SessionMap = make(map[int64]*UserSession)

func WsHandler(c *gin.Context) {
	tmp, exist := c.Get("user")
	ecustUser := tmp.(*user.EcustUser)
	if !exist {
		c.String(http.StatusUnauthorized, "not login")
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)

	session := &UserSession{
		Conn:        conn,
		User:        ecustUser,
		SendChannel: make(chan *SendingMessage, 100),
	}
	SessionMap[ecustUser.UserId] = session
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to upgrade")
		return
	}

	// 发送channel
	util.SafeGo(func() {
		for {
			sendingMessage := <-session.SendChannel
			_ = conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
			if err := conn.WriteMessage(sendingMessage.MessageType, sendingMessage.Data); err != nil {
				delete(SessionMap, ecustUser.UserId)
				_ = conn.Close()
				break
			}
		}
	})

	// 接收channel
	util.SafeGo(func() {
		for {
			messageType, bytes, err := conn.ReadMessage()
			if err != nil {
				delete(SessionMap, ecustUser.UserId)
				_ = conn.Close()
				break
			}
			var packet = &dto.Packet{}
			if messageType == websocket.BinaryMessage {
				// protobuf
				if err := packet.Unmarshal(bytes); err != nil {
					log.Warnf("failed to unmarshal binary bytes, %+v", err)
					continue
				}
			} else if messageType == websocket.TextMessage {
				// json
				if err := json.Unmarshal(bytes, packet); err != nil {
					log.Warnf("failed to unmarshal text, %+v", err)
					continue
				}
			}

			// TODO 加群/加好友请求 群/好友消息处理
		}
	})
}
