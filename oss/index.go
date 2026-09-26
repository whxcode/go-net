// 处理从前端接受文件相关操作
package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"go-net/config"
	"go-net/pool"
	"go-net/std"
)

type FileMeta struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}

func StorageFile(filename string, fd multipart.File) string {
	file, err := createFile(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	io.Copy(file, fd)

	return filename
}

func StorageMetaFile(filename string, data []byte) string {
	file, err := CreateMetaFile(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	io.Copy(file, bytes.NewReader(data))

	return filename
}

func StorageFiles(fileHeaderies []*multipart.FileHeader) []string {
	hashWorker := pool.NewWorkers(len(fileHeaderies))
	type FileObj struct {
		file multipart.File
		meta *multipart.FileHeader
	}

	fileMap := std.NewHash[FileObj]()

	result := std.NewQueue[string]()

	for _, head := range fileHeaderies {
		v, _ := head.Open()
		defer v.Close()

		hashWorker.Post(func(args ...interface{}) {
			hash := getHashWithFile(v)
			result.Push(hash)

			if !IsHashExits(hash) {
				fmt.Println("--- 文件已经存在---")
				return
			}

			// 2. 回到开头
			if _, err := v.Seek(0, io.SeekStart); err != nil {
				panic(err)
			}

			fileMap.Set(hash, FileObj{
				file: v,
				meta: args[0].(*multipart.FileHeader),
			})
		}, head)
	}

	hashWorker.Wait()

	if !fileMap.IsEmpty() {

		storageWorker := pool.NewWorkers(len(fileHeaderies))

		for key, obj := range fileMap.Raw() {
			storageWorker.Post(func(args ...interface{}) {
				StorageFile(key, obj.file)

				data, _ := json.Marshal(FileMeta{
					Filename: obj.meta.Filename,
					Size:     obj.meta.Size,
					Hash:     key,
				})

				StorageMetaFile(key, data)
			})
		}

		storageWorker.Wait()

	}

	return result.Raw()
}

const tsecretKey = "your_secret_key_here"

func GetFile(file string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%s/%s", config.ConfigData.Server.FileOss, file))
}

func Signature(message string) string {
	h := hmac.New(sha256.New, []byte(tsecretKey))
	h.Write([]byte(message))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// 根据 hash 生成一个带时效的临时下载地址
//
// 地址指向下载路由 /api/file/download/:hash（常量 controller.KFileDownload），
// 路径里带 expred + signature，由 DownloadMiddleware 校验。
// 注意：这里拼的是本机地址，部署到服务器/域名后要改成对外可访问的地址（域名或 IP:端口）。
func GenerateSignedURL(hash string, expred time.Duration) string {
	expredAt := time.Now().Add(expred).Unix()
	expredAtStr := strconv.FormatInt(expredAt, 10)

	message := fmt.Sprintf("%s:%s", hash, expredAtStr)

	signature := Signature(message)

	baseURL := fmt.Sprintf("http://localhost%s/api/file/download", config.ConfigData.Server.Port)

	signatureURL := fmt.Sprintf("%s/%s?expred=%s&signature=%s", baseURL, hash, expredAtStr, signature)

	return signatureURL
}
