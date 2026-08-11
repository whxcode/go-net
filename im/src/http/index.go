package httpServer

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"go-net/im/src/http/controll"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 自定义响应中间件
func responseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 执行后续 handler

		// 获取 handler 返回的数据（通过 c.Set 传递）
		data, exists := c.Get("response")
		if !exists {
			return
		}

		// 统一格式化
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": data,
		})
	}
}

// 定义 接口
func resetApi(r *gin.Engine) {
	r.POST("/news", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.PUT("/edit", func(c *gin.Context) {
		c.String(200, "编辑数据 Hello World")
	})

	r.DELETE("/delete", func(c *gin.Context) {
		c.String(200, "删除 Hello World")
	})
}

// 处理返回数据
func responseData(r *gin.Engine) {
	type Person struct {
		Title string `json:"title"`
		Age   uint8  `json:"age"`
		Dogs  []int
	}

	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, &Person{Title: "王恒星"})
	})

	r.GET("/jsonp", func(c *gin.Context) {
		c.JSONP(http.StatusOK, &Person{Title: "王恒星"})
	})

	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, &Person{Title: "王恒星"})
	})

	r.GET("/html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", Person{Title: "王恒星", Age: 10, Dogs: []int{1, 2, 3, 4}})
	})

	/*
		r.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]string{
				"name": "王恒星",
			})
			// c.HTML(200, "<html><body><h1>Data</h1></body></html>", "whx")
			// c.String(200, "Hello World")
		})
	*/
}

func getRequestParams(r *gin.Engine) {
	r.GET("/d1", func(c *gin.Context) {
		name := c.Query("name")
		id := c.Query("id")

		c.String(200, "Hello %s : %s", id, name)
	})

	r.POST("/d2", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
			Pwd  string `json:"pwd"`
		}

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
	})

	r.POST("/d3", func(c *gin.Context) {
		name := c.PostForm("name")
		pwd := c.PostForm("pwd")
		file, err := c.FormFile("file") // "file" 是前端表单的字段名
		if err != nil {
		}

		c.JSON(http.StatusOK, gin.H{
			"name":     name,
			"pwd":      pwd,
			"file":     file.Filename,
			"fileSize": file.Size,
		})
	})

	r.POST("/d4", func(c *gin.Context) {
		var req struct {
			Name string                `form:"name" json:"name"`
			Pwd  string                `form:"pwd" json:"pwd"`
			File *multipart.FileHeader `form:"file" json:"file"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filename := "1"
		fileSize := 0

		if req.File != nil {
			filename = req.File.Filename
			fileSize = int(req.File.Size)
		}

		c.JSON(http.StatusOK, gin.H{
			"name":     req.Name,
			"pwd":      req.Pwd,
			"file":     filename,
			"fileSize": fileSize,
		})
	})

	r.PUT("/edit", func(c *gin.Context) {
		c.String(200, "编辑数据 Hello World")
	})

	r.DELETE("/delete", func(c *gin.Context) {
		c.String(200, "删除 Hello World")
	})
}

func configLog() io.Writer {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 以追加的方式
	// 以追加的方式
	f, _ := os.OpenFile("./log/gin.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	logrus.SetOutput(io.MultiWriter(f, os.Stdout))

	return logrus.StandardLogger().Out
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// 记录到 logrus（自动 JSON）
		logrus.WithFields(logrus.Fields{
			"time":      start.Format("2006-01-02 15:04:05"),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    c.Writer.Status(),
			"latency":   latency.String(),
			"clientIP":  c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}).Info("HTTP Request")
	}
}

func group(_r *gin.Engine) {
	r := _r.Group("/user")
	{
		r.GET("/d1", func(c *gin.Context) {
			c.String(200, "Hello World User d1")
		})
	}
}

func groupControll(_r *gin.Engine) {
	r := _r.Group("/controll")
	{
		r.GET("/d1", func(c *gin.Context) {
			c.String(200, "Hello World User d1")
		})
	}
}

func Start() {
	configLog()
	// gin.DefaultWriter = configLog()
	r := gin.New()
	r.Use(LoggerMiddleware()) // 自定义日志中间件
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("/home/whx/study/go-net/im/pages/*.html")
	// 使用默认CORS中间件，允许所有跨域请求
	r.Use(cors.Default())

	r.Use(responseMiddleware()) // 使用自定义响应中间件

	group(r)

	for k, v := range controll.GetControllMap {
		r.GET(k, func(c *gin.Context) {
			// 调用 handler 并获取返回值
			result := v(c)
			// 将返回值设置到上下文中
			c.Set("response", result)
		})
	}

	for k, v := range controll.PostControllMap {
		r.POST(k, func(c *gin.Context) {
			// 调用 handler 并获取返回值
			result := v(c)
			// 将返回值设置到上下文中
			c.Set("response", result)
		})
	}

	r.Static("/asset", "/home/whx/study/go-net/im/pages/asset/")

	responseData(r)
	getRequestParams(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
