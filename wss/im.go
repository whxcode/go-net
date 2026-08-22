package wss

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"go-net/middleware"
	"go-net/model"
	"go-net/pool"
	"go-net/redis"
)

func RequestWsHandle(c *gin.Context) middleware.CloseHandle {
	wsConn, ok := c.Get("wsConn")

	//conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !ok {
		fmt.Println("wsConn not found in context")
		return nil
	}

	conn, o := wsConn.(*websocket.Conn)

	if !o {
		fmt.Println("wsConn is not of type *websocket.Conn")
		return nil
	}

	var userID uint = 0

	close := func(conn *websocket.Conn) {
		// 用户主动关闭链接

		pool.UserPool.RemoveUser(userID)
		conn.Close()
	}

	// 客户端首次链接
	userIDStr := c.Query("userID")

	userIDUint, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID = uint(userIDUint)

	pool.UserPool.AddUser(userID, conn)

	// 每当用户上上线后；推送离线信息
	messages, _ := redis.GetOfflineMessage(userID)

	if messages != nil {
		for _, msg := range messages {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	return close
}

func ChannelHandle(conn *websocket.Conn, p []byte, message *model.Message) bool {
	senderConn := pool.UserPool.GetUserConn(message.ReceiverID)

	// 保存数据
	go func() { model.MessageDB.Save(message) }()

	// 将离线信息存入 redis 7 天时间。
	if senderConn == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("当前好友不在线"))
		redis.SaveOfflineMessage(message.ReceiverID, p)
		return true
	}

	senderConn.WriteMessage(websocket.TextMessage, p)
	return true
}
