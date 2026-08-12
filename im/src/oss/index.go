// 处理从前端接受文件相关操作
package oss

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	config "go-net/im/src"
	"go-net/im/src/utils"
)

func Init() {
}

func makeFile(filename string) (*os.File, error) {
	dir := fmt.Sprintf("%s/%s", config.ConfigData.Server.FileOss, utils.GetToday())

	os.MkdirAll(dir, 0o755)

	return os.Create(fmt.Sprintf("%s/%s", dir, filename))
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
