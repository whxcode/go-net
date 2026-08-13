package middleware

import "github.com/gin-gonic/gin"

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
