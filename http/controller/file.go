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

type UploadResponse struct {
	// 文件 hash 地址
	Hash string `json:"hash"`
	// 文件名称
	Filename string `json:"filename"`
	// 文件大小 字节
	Size int64 `json:"size"`
}

// @Summary  保存用户上传的文件
// @Tags 文件
// @Accept multipart/form-data
// @Param files formData []file true "文件列表"
// @Success 200 {object} utils.KResponse{data=[]UploadResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /file/upload [post]
func (*fileController) Upload(c *gin.Context) *utils.KResponse {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := form.File["files"]

	result := []*UploadResponse{}

	for _, v := range files {
		result = append(result, &UploadResponse{
			Hash:     oss.StorageFile(v),
			Filename: v.Filename,
			Size:     v.Size,
		})
	}

	return utils.MakeResponse(result)
}

type GetFileRequest struct {
	Files []string `json:"files"`
}

type GetFileResponse struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

// @Summary  根据 hash 获取文件的 url 地址
// @Description 此接口生成的 url 地址；含有时效性；不建议作为长时间显示。
// @Tags 文件
// @Param request body GetFileRequest true "文件的hash列表"
// @Success 200 {object} utils.KResponse{data=[]GetFileResponse} "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /file/getfile [post]
func (*fileController) GetFile(c *gin.Context) *utils.KResponse {
	parmas := &GetFileRequest{}

	err := c.ShouldBindJSON(parmas)
	if err != nil {
	}

	// os.Open()

	result := make([]*GetFileResponse, 0, len(parmas.Files))

	for _, hash := range parmas.Files {
		r := &GetFileResponse{
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

// @Summary  根据 url 地址；下载文件;
// @Tags 文件
// @Param hash path string true "哈希地址"
// @Param expred query string true "过期时间"
// @Param signature query string true "签名"
// @Success 200 {object} []byte "成功"
// @Failure 500 {object} utils.KResponse "服务器错误"
// @Router /download/:hash [get]
func (*fileController) DowloadFile(c *gin.Context) {
	hash := c.Param("hash")

	c.File(config.ConfigData.Server.FileOss + "/" + hash)
}
