package controll

import (
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
