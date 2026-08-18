package controller

import (
	"strconv"
	"time"

	"go-net/config"
	"go-net/oss"
	"go-net/utils"

	"github.com/gin-gonic/gin"
)

type fileController struct{}

var FileController = &fileController{}

/**
* @beif 保存用户上传的文件
* @parma files os.File 文件列表
* @return []
* **/
func (*fileController) Upload(c *gin.Context) *utils.KResponse {
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

	return utils.MakeResponse(result)
}

func (*fileController) GetFile(c *gin.Context) *utils.KResponse {
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

	return utils.MakeResponse(result)
}

func (*fileController) DownloadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

func (*fileController) DowloadFile(c *gin.Context) {
	hash := c.Param("hash")

	c.File(config.ConfigData.Server.FileOss + "/" + hash)
}
