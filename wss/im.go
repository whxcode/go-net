package wss

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	dbGroup "go-net/db/group"
	dbMessage "go-net/db/message"
	"go-net/model"
	"go-net/pool"
	MessageRedis "go-net/redis/message"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
)

func broadcastFriendMessage(message *model.Message, msg []byte) {
	senderConn := pool.UserPool.GetUserConn(message.ReceiverID)

	// 用户在线
	if senderConn != nil {
		senderConn.WriteMessage(websocket.TextMessage, msg)
	}
}

func broadcastGroupMessage(message *model.Message, msg []byte) {
	result := dbGroup.GroupDB.GroupMembers(message.ReceiverID)

	for _, userID := range result {
		if userID == message.SenderID {
			continue
		}

		senderConn := pool.UserPool.GetUserConn(userID)
		if senderConn != nil {
			senderConn.WriteMessage(websocket.TextMessage, msg)
		}

	}
}

func IM(c *gin.Context) {
	conn, err := pool.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("升级为WebSocket失败:", err)
		return
	}

	defer conn.Close()

	fmt.Println("WebSocket连接已建立:", conn.RemoteAddr())

	conn.SetReadDeadline(time.Now().Add(PongWait))

	fmt.Println(" conn.SetReadDeadline(time.Now().Add(PongWait)) ")

	// 客户端首次链接
	userIDStr := c.Query("userID")
	userIDUint, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userIDUint)
	fmt.Println("---", userID, "---", userIDUint, "---", userIDStr)

	pool.UserPool.AddUser(uint(userID), conn)

	offlineMsg := MessageRedis.GetOfflineMessage(c, userID)

	// 推送离线消息
	if len(offlineMsg) > 0 {
		for _, msg := range offlineMsg {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	for {
		fmt.Println("-------------------------------------- start -------------------------------")
		_, msg, err := conn.ReadMessage()
		if err != nil {

			// 判断是否超时
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("读取消息超时，关闭连接")
				break
			}

			// 判断是否是客户端主动关闭
			if closeErr, ok := err.(*websocket.CloseError); ok {
				fmt.Printf("客户端关闭连接，状态码: %d，原因: %s\n", closeErr.Code, closeErr.Text)
				break
			}

			fmt.Println("--->err:", err)
			panic(err)

		}

		if len(msg) == 0 {
			fmt.Println(userID, "收到空消息---")
			continue
		}

		conn.SetReadDeadline(time.Now().Add(PongWait))

		var message *model.Message
		err = json.Unmarshal(msg, &message)

		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("解析JSON失败", err)
			break
		}

		switch message.Type {
		case model.ChannelTypePING:
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"type":%v}`, model.ChannelTypePONG))) // PONG
		case model.ChannelTypeFriend:
			broadcastFriendMessage(message, msg)
		case model.ChannelTypeGroup:
			broadcastGroupMessage(message, msg)
		}

		switch message.Type {
		case model.ChannelTypeFriend, model.ChannelTypeGroup:
			go func() { dbMessage.MessageDB.Save(message) }()
		}

		fmt.Println("-------------------------------------- end -------------------------------")

	}
}
