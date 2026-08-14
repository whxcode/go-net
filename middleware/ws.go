package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"go-net/model"
	"go-net/pool"
	"go-net/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func BeatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := pool.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		c.Set("wsConn", conn)

		c.Next()
	}
}

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
)

type (
	CloseHandle     = func(conn *websocket.Conn)
	RequestWsHandle = func(c *gin.Context) CloseHandle
	ChannelHandle   = func(conn *websocket.Conn, p []byte, message *model.Message) bool
)

func BeatExecute(
	requestHandle RequestWsHandle,
	channelHandle ChannelHandle,
) func(c *gin.Context) {
	return func(c *gin.Context) {
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

		// var userID uint = 0
		closeHandle := requestHandle(c)

		if closeHandle == nil {
			conn.Close()
			return
		}

		defer closeHandle(conn)

		conn.SetReadDeadline(time.Now().Add(PongWait))

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("ReadMessage err:", err)
				break
			}
			// 收到任何消息重置读超时
			conn.SetReadDeadline(time.Now().Add(PongWait))

			message := &model.Message{}

			if err := json.Unmarshal(msg, &message); err != nil {
				println("json.Unmarshal err:", err)
				return
			}

			fmt.Println("message:[", message.Type, "]")

			if message.Type == model.ChannelTypePING {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type":2}`)) // PONG
				continue
			}

			message.MsgID = utils.GenerateSnowflakeID()

			if channelHandle(conn, msg, message) {
			} else {
				return
			}

		}
	}
}
