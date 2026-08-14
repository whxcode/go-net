package wss

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"go-net/model"
	"go-net/pool"

	"github.com/gin-gonic/gin"
)

func HandleWSS(c *gin.Context) {
	wsConn, ok := c.Get("wsConn")

	//conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !ok {
		fmt.Println("wsConn not found in context")
		return
	}

	conn, o := wsConn.(*websocket.Conn)

	if !o {
		fmt.Println("wsConn is not of type *websocket.Conn")
		return
	}

	var userID uint = 0

	defer func() {
		// 用户主动关闭链接

		if userID == 0 {
			pool.UserPool.RemoveUser(userID)
		}

		conn.Close()
	}()

	timeout := 5 * time.Second

	// 客户端首次链接

	userIDStr := c.Query("userID")

	userIDUint, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID = uint(userIDUint)

	pool.UserPool.AddUser(userID, conn)

	conn.SetReadDeadline(time.Now().Add(timeout))

	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(timeout))

		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(timeout))
	})

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(timeout))
		return nil
	})

	conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

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
