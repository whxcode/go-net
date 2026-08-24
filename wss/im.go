package wss

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

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
		fmt.Println("--关闭链接--")
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

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
)

func IM(c *gin.Context) {
	conn, err := pool.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("升级为WebSocket失败:", err)
		return
	}

	defer func() {
		fmt.Println("WebSocket连接已关闭:", conn.RemoteAddr())
		conn.Close()
	}()

	fmt.Println("WebSocket连接已建立:", conn.RemoteAddr())

	conn.SetReadDeadline(time.Now().Add(PongWait))

	fmt.Println(" conn.SetReadDeadline(time.Now().Add(PongWait)) ")

	// 客户端首次链接
	userIDStr := c.Query("userID")
	userIDUint, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userIDUint)
	fmt.Println("---", userID, "---", userIDUint, "---", userIDStr)

	pool.UserPool.AddUser(uint(userID), conn)

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

		if message.Type == model.ChannelTypePING {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type":2}`)) // PONG
			continue
		}

		senderConn := pool.UserPool.GetUserConn(message.ReceiverID)

		if senderConn != nil {
			senderConn.WriteMessage(websocket.TextMessage, msg)
		}

		fmt.Println("-------------------------------------- end -------------------------------")

	}
}
