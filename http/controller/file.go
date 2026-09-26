package controller

import (
	"strconv"
	"time"

	"go-net/oss"
	"go-net/response"

	"github.com/gin-gonic/gin"
)

type fileController struct{}

var FileController = &fileController{}

// Upload 保存用户上传的文件
//
// 每个文件按内容 sha256 作为 hash 落盘（同时写一份 <hash>meta.json 记录文件名和大小），
// 同样的内容重复上传只会命中已有的 hash，不会重复存储。
// 返回值为文件 hash 列表；前端用 hash 拼预览地址 GET /api/file/preview/{hash} 展示。
//
// @Summary  上传文件（可多选）
// @Tags 文件
// @Accept multipart/form-data
// @Param files formData []file true "文件列表"
// @Success 200 {object} response.KResponse{data=[]string} "成功，data 是文件 hash 列表"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /file/upload [post]
func (*fileController) Upload(c *gin.Context) *response.KResponse {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	files := form.File["files"]
	result := oss.StorageFiles(files)

	return response.MakeResponse(result)
}

// SignURLsRequest 批量换取临时下载地址的请求体
type SignURLsRequest struct {
	Files []string `json:"files"`
}

// SignURLsResponse 单个 hash 对应的临时下载地址
type SignURLsResponse struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

// SignURLs 批量把文件 hash 换成带时效的临时下载地址
//
// 生成地址形如 /api/file/download/{hash}?expred=<过期unix秒>&signature=<签名>，
// 由 DownloadMiddleware 校验，默认 1 小时有效；一般用于需要分享/外链的场景。
// 站内展示图片、头像请直接用 PreviewFile，不要走这里。
//
// @Summary  批量获取文件的临时下载地址
// @Description 返回的 url 含时效性（默认 1 小时），过期或签名不符会返回 403；不建议作为长期显示地址。
// @Tags 文件
// @Param request body SignURLsRequest true "文件的 hash 列表"
// @Success 200 {object} response.KResponse{data=[]SignURLsResponse} "成功"
// @Failure 500 {object} response.KResponse "服务器错误"
// @Router /file/signurls [post]
func (*fileController) SignURLs(c *gin.Context) *response.KResponse {
	parmas := &SignURLsRequest{}

	err := c.ShouldBindJSON(parmas)
	if err != nil {
	}

	result := make([]*SignURLsResponse, 0, len(parmas.Files))

	for _, hash := range parmas.Files {
		r := &SignURLsResponse{
			Hash: hash,
			Url:  oss.GenerateSignedURL(hash, time.Hour),
		}

		result = append(result, r)

	}

	return response.MakeResponse(result)
}

// PreviewFile 根据 hash 直接返回文件内容
// @Summary 预览文件（按 hash 直出）
// @Description 注意该接口没有做任何权限验证；慎用；只给站内展示用，需要对外分享请走 SignURLs 换临时地址。
// @Tags 文件
// @Produce application/octet-stream
// @Param hash path string true "文件 hash"
// @Success 200 {file} binary "文件二进制数据"
// @Failure 404 {object} response.KResponse "文件不存在"
// @Router /file/preview/{hash} [get]
func (*fileController) PreviewFile(c *gin.Context) {
	hash := c.Param("hash")

	_, filepath := oss.MakeOssStorageFilePath(hash)

	c.File(filepath)
}

// PreviewMetaFile 根据 hash 直接返回文件元信息（落盘时写入的 <hash>meta.json）
// @Summary 预览文件元信息（JSON）
// @Description 返回 filename/size/hash 三个字段；无鉴权，主要给开发排查用；前端展示不需要调它。
// @Tags 文件
// @Produce application/json
// @Param hash path string true "文件 hash"
// @Success 200 {object} oss.FileMeta "文件元信息"
// @Failure 404 {object} response.KResponse "元信息不存在"
// @Router /file/preview/{hash}/meta [get]
func (*fileController) PreviewMetaFile(c *gin.Context) {
	hash := c.Param("hash")

	_, filepath := oss.MakeOssStorageMetaFilePath(hash)

	c.File(filepath)
}

// DownloadMiddleware 临时下载地址的签名校验中间件
//
// 只挂在 KFileDownload（/download/:hash）这一条路由上，校验顺序：
//  1. expred 必须是合法的 unix 秒时间戳，否则 400；
//  2. 当前时间不能超过 expred，否则 403（链接过期）；
//  3. signature 必须等于 HMAC(hash:expred)，否则 403。
//
// 失败响应是 {"error": "..."}，不走 KResponse 统一包装（这里是下载链路，不是业务接口）。
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

// DownloadFile 校验签名后下载文件
// @Summary 通过临时签名地址下载文件
// @Description 需要 SignURLs 接口生成的 expred + signature，缺一不可；校验失败返回 400/403，body 是 {"error":"..."}。
// @Tags 文件
// @Produce application/octet-stream
// @Param hash path string true "文件 hash"
// @Param expred query string true "过期时间（unix 秒）"
// @Param signature query string true "签名"
// @Success 200 {file} binary "文件二进制数据"
// @Failure 404 {object} response.KResponse "文件不存在"
// @Router /file/download/{hash} [get]
func (*fileController) DownloadFile(c *gin.Context) {
	hash := c.Param("hash")
	_, filepath := oss.MakeOssStorageFilePath(hash)

	c.File(filepath)
}
