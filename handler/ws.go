package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/lz1998/ecust_im/dto"
	"github.com/lz1998/ecust_im/model"
	"github.com/lz1998/ecust_im/model/friend"
	"github.com/lz1998/ecust_im/model/group"
	"github.com/lz1998/ecust_im/model/group_member"
	"github.com/lz1998/ecust_im/model/request"
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
	Conn   *websocket.Conn
	User   *user.EcustUser
	PubSub *redis.PubSub
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
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to upgrade")
		return
	}

	pubSub := model.RDb.Subscribe(context.Background(), fmt.Sprintf("PUSH:%d", ecustUser.UserId))

	session := &UserSession{
		Conn:   conn,
		User:   ecustUser,
		PubSub: pubSub,
	}
	SessionMap[ecustUser.UserId] = session

	// 发送channel
	util.SafeGo(func() {
		for {
			streamName := fmt.Sprintf("PACKET:%d", ecustUser.UserId)
			// TODO 这个可以放在register时创建
			model.RDb.XGroupCreate(context.Background(), streamName, "cg", "0-0")
			xStream, err := model.RDb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Streams:  []string{streamName},
				Group:    "cg",
				Consumer: "c",
				Count:    1,
				Block:    1 * time.Second,
				NoAck:    false,
			}).Result()
			if err != nil {
				log.Warnf("read redis queue error")
				continue
			}
			for _, stream := range xStream {
				for _, message := range stream.Messages {
					pid, ok := message.Values["packetId"]
					if !ok {
						log.Warnf("packetId not exists")
						model.RDb.XAck(context.Background(), streamName, "cg1", message.ID)
						continue
					}
					packetId, ok := pid.(string)
					if !ok {
						log.Warnf("packet is not string")
						model.RDb.XAck(context.Background(), streamName, "cg1", message.ID)
						continue
					}
					data, err := model.LDb.Get([]byte(packetId), nil)
					if err != nil {
						log.Warnf("failed to get packet, err: %+v", err)
						continue
					}
					if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
						delete(SessionMap, ecustUser.UserId)
						_ = pubSub.Close()
						_ = conn.Close()
						break
					}
				}
			}
		}
	})

	// 接收channel
	util.SafeGo(func() {
		for {
			messageType, bytes, err := conn.ReadMessage()
			if err != nil {
				delete(SessionMap, ecustUser.UserId)
				_ = pubSub.Close()
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

			util.SafeGo(func() {
				HandlePacket(ecustUser.UserId, packet)
			})
		}
	})
}

// 向客户端发送数据包
// leveldb: {msg_id1: bytes1, msg_id2: bytes2, msg_id3: bytes3}
// redis: [msg_id1, msg_id2, msg_id3]
func SendPacket(userId int64, packet *dto.Packet) error {
	packetData, err := packet.Marshal()
	if err != nil {
		return err
	}
	// 放到 leveldb数据库 和 redis队列
	packetId := GeneratePacketId(packet)

	// 保存leveldb
	if err = model.LDb.Put([]byte(packetId), packetData, nil); err != nil {
		return err
	}

	streamName := fmt.Sprintf("PACKET:%d", userId)
	// redis 消息队列发布
	return model.RDb.XAdd(context.Background(), &redis.XAddArgs{
		ID:     "*",
		Stream: streamName,
		Values: []string{"packetId", packetId},
	}).Err()
}

func GeneratePacketId(packet *dto.Packet) string {
	// request:friend:123
	// request:group:123
	// msg:friend:<userId>:123
	// msg:group:<groupId>:123
	uniqId := util.GenerateId()
	var packetId string
	if packet.PacketType == dto.Packet_TMsg {
		msg := packet.GetMsg()
		if msg.MsgType == dto.Msg_TFriend {
			packetId = fmt.Sprintf("msg:friend:%d:%d", msg.ToId, uniqId)
		} else {
			packetId = fmt.Sprintf("msg:group:%d:%d", msg.ToId, uniqId)
		}
	} else {
		req := packet.GetRequest()
		if req.ReqType == dto.Request_TFriend {
			packetId = fmt.Sprintf("request:friend:%d:%d", req.ReqId, uniqId)
		} else {
			packetId = fmt.Sprintf("request:group:%d:%d", req.ReqId, uniqId)
		}
	}
	return packetId
}

func HandlePacket(fromUserId int64, packet *dto.Packet) {
	if packet.PacketType == dto.Packet_TMsg {
		msg := packet.GetMsg()
		msg.FromId = fromUserId
		HandleMsg(msg)
	} else {
		req := packet.GetRequest()
		req.FromId = fromUserId
		HandleRequest(req)
	}
}

func HandleRequest(req *dto.Request) {
	// 保存MySQL
	// 转发给对方/群主
	r := &request.EcustRequest{
		ReqType: int64(req.ReqType),
		FromId:  req.FromId,
		ToId:    req.ToId,
	}
	modelRequest, err := request.CreateRequest(r)
	if err != nil {
		log.Warnf("failed to create request, err: %+v", err)
		return
	}
	req.ReqId = modelRequest.ReqId
	packet := &dto.Packet{
		Timestamp:  time.Now().Unix(),
		PacketType: dto.Packet_TRequest,
		Data: &dto.Packet_Request{
			Request: req,
		},
	}
	if req.ReqType == dto.Request_TFriend {
		// 发送给对方
		if err := SendPacket(req.ToId, packet); err != nil {
			log.Errorf("failed to send packet, err: %+v", err)
			return
		}
	} else {
		// 发送给群主
		ecustGroup, err := group.GetGroup(req.ToId)
		if err != nil {
			log.Errorf("failed to get group, err: %+v", err)
			return
		}
		if err := SendPacket(ecustGroup.OwnerId, packet); err != nil {
			log.Errorf("failed to send packet, err: %+v", err)
			return
		}
	}

}

func HandleMsg(msg *dto.Msg) {
	// 私聊 直接转发给对方
	// 群聊 找出群内所有人，发送到每个人的队列（写扩散）
	packet := &dto.Packet{
		Timestamp:  time.Now().Unix(),
		PacketType: dto.Packet_TMsg,
		Data: &dto.Packet_Msg{
			Msg: msg,
		},
	}
	if msg.MsgType == dto.Msg_TFriend {
		if !friend.IsFriend(msg.FromId, msg.ToId) {
			log.Warnf("not friend: %d %d", msg.FromId, msg.ToId)
			return
		}
		// 发送给对方
		if err := SendPacket(msg.ToId, packet); err != nil {
			log.Errorf("failed to send packet, err: %+v", err)
			return
		}
	} else {
		if !group_member.IsInGroup(msg.ToId, msg.FromId) {
			log.Warnf("not in group, %d %d", msg.ToId, msg.FromId)
			return
		}

		// 发送给每个人
		groupId := msg.ToId
		memberIds, err := group_member.ListGroupMember(groupId)
		if err != nil {
			log.Errorf("failed to list group member")
			return
		}
		for _, memberId := range memberIds {
			if err := SendPacket(memberId, packet); err != nil {
				log.Errorf("failed to send packet, err: %+v", err)
				return
			}
		}
	}
}
