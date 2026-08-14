package wss

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"go-net/middleware"
	"go-net/model"
	"go-net/pool"
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

	return close
}

func ChannelHandle(conn *websocket.Conn, p []byte, message *model.Message) bool {
	senderConn := pool.UserPool.GetUserConn(message.ReceiverID)

	if senderConn == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("当前好友不在线"))
		return true
	}

	senderConn.WriteMessage(websocket.TextMessage, p)
	return true
}
