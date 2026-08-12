package controll

import (
	"strconv"
	"strings"
	"time"

	config "go-net/im/src"
	"go-net/im/src/oss"

	"github.com/gin-gonic/gin"
)

/**
* 保存用户上传的文件
* **/
func Upload(c *gin.Context) *KResponse {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := form.File["files"]

	type P struct {
		Hash     string `json:"hash"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}

	result := []*P{}

	for _, v := range files {
		result = append(result, &P{
			Hash:     oss.StorageFile(v),
			Filename: v.Filename,
			Size:     v.Size,
		})
	}

	return MakeResponse(result)
}

func GetFile(c *gin.Context) *KResponse {
	type P struct {
		Files []string `json:"files"`
	}

	parmas := &P{}

	err := c.ShouldBindJSON(parmas)
	if err != nil {
	}

	// os.Open()

	type R struct {
		Url  string `json:"url"`
		Hash string `json:"hash"`
	}
	result := make([]*R, 0, len(parmas.Files))

	for _, hash := range parmas.Files {
		r := &R{
			Hash: hash,
			Url:  oss.GenerateSignedURL(hash, time.Hour),
		}

		result = append(result, r)

	}

	return MakeResponse(result)
}

func DownloadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if !strings.HasPrefix(path, "/file/download/") {
			c.Next()
			return
		}

		expred := c.Query("expred")
		expredInt, err := strconv.ParseInt(expred, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid expred parameter"})
			return
		}

		if time.Now().Unix() > expredInt {
			c.AbortWithStatusJSON(403, gin.H{"error": "link expired"})
			return
		}

		signature := c.Query("signature")
		newS := oss.Signature(c.Param("hash") + ":" + expred)

		if newS != signature {
			c.AbortWithStatusJSON(403, gin.H{"error": "invalid signature:" + newS})
			return
		}

		c.Next()
	}
}

func DowloadFile(c *gin.Context) {
	hash := c.Param("hash")

	c.File(config.ConfigData.Server.FileOss + "/" + hash)
}
