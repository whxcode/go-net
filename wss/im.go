package wss

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-net/model"
	"go-net/pool"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWSS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	var userID uint = 0

	defer func() {
		// 用户主动关闭链接

		if userID == 0 {
			pool.UserPool.RemoveUser(userID)
		}

		conn.Close()
	}()

	// 客户端首次链接

	userIDStr := c.Query("userID")

	userIDUint, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID = uint(userIDUint)

	pool.UserPool.AddUser(userID, conn)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		message := &model.Message{}

		if err := json.Unmarshal(msg, &message); err != nil {
			println("json.Unmarshal err:", err)
			return
		}

		senderConn := pool.UserPool.GetUserConn(message.ReceiverID)

		if senderConn == nil {
			conn.WriteMessage(msgType, []byte("当前好友不在线"))
			continue
		}

		senderConn.WriteMessage(msgType, msg)
	}
}
