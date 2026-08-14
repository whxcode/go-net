package middleware

import (
	"go-net/pool"

	"github.com/gin-gonic/gin"
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
