package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var errMsg string
				var funcName string
				var file string
				var line int

				// 遍历调用栈，找到第一个非 runtime/gin 的帧
				pcs := make([]uintptr, 20)
				n := runtime.Callers(2, pcs)
				frames := runtime.CallersFrames(pcs[:n])

				for {
					frame, more := frames.Next()

					// 跳过 runtime 和 gin 框架
					if strings.Contains(frame.File, "runtime/") ||
						strings.Contains(frame.File, "gin@") ||
						strings.Contains(frame.File, "net/http") {
						if !more {
							break
						}
						continue
					}

					// 找到你自己的代码
					funcName = frame.Function
					file = frame.File
					line = frame.Line

					// 只保留文件名
					for i := len(file) - 1; i >= 0; i-- {
						if file[i] == '/' {
							file = file[i+1:]
							break
						}
					}
					break
				}

				switch v := r.(type) {
				case error:
					if errors.Is(v, gorm.ErrRecordNotFound) {
						c.AbortWithStatusJSON(http.StatusOK, gin.H{
							"code":    http.StatusNotFound,
							"message": http.StatusText(http.StatusNotFound),
						})
						return
					}
					errMsg = v.Error()
				case string:
					errMsg = v
				default:
					errMsg = fmt.Sprintf("%v", v)
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf("[%s:%d] %s", file, line, errMsg),
					"func":    funcName,
				})
			}
		}()

		c.Next()
	}
}
