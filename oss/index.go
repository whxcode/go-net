// 处理从前端接受文件相关操作
package oss

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	config "go-net/im/src"
)

func Init() {
}

func makeFile(filename string) (*os.File, error) {
	file := fmt.Sprintf("%s/%s", config.ConfigData.Server.FileOss, filename)
	return os.Create(file)
}

func HashBytes(data []byte) string {
	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}

func HashFile(data multipart.File) string {
	hash := sha256.New()

	if _, err := io.Copy(hash, data); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func StorageFile(f *multipart.FileHeader) string {
	fd, err := f.Open()
	if err != nil {
		panic(err)
	}

	defer fd.Close()

	filename := HashFile(fd)

	// 2. 回到开头
	if _, err := fd.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	file, err := makeFile(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	io.Copy(file, fd)

	return filename
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

// 根据 hash 生成一个临时的url 地址
func GenerateSignedURL(hash string, expred time.Duration) string {
	expredAt := time.Now().Add(expred).Unix()
	expredAtStr := strconv.FormatInt(expredAt, 10)

	message := fmt.Sprintf("%s:%s", hash, expredAtStr)

	signature := Signature(message)

	baseURL := "http://localhost:8080/file/download"

	signatureURL := fmt.Sprintf("%s/%s?expred=%s&signature=%s", baseURL, hash, expredAtStr, signature)

	return signatureURL
}
