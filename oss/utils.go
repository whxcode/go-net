package oss

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"go-net/config"
)

func GetHashWithBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func getHashWithFile(data multipart.File) string {
	hash := sha256.New()

	if _, err := io.Copy(hash, data); err != nil {
		panic(err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

/**
* hash 文件的hash；也是文件存储地址。
*
* */
func MakeOssStorageFilePath(hash string) (dir string, filepath string) {
	dir = fmt.Sprintf("%s/%s", config.ConfigData.Server.FileOss, hash[:2])
	filepath = fmt.Sprintf("%s/%s", dir, hash)

	return
}

/**
* hash 文件的hash；也是文件存储地址。
*
* */
func MakeOssStorageMetaFilePath(hash string) (dir string, filepath string) {
	dir = fmt.Sprintf("%s/%s", config.ConfigData.Server.FileOss, hash[:2])
	filepath = fmt.Sprintf("%s/%smeta.json", dir, hash)

	return
}

func IsHashExits(hash string) bool {
	_, filepath := MakeOssStorageMetaFilePath(hash)
	_, err := os.Stat(filepath)

	return os.IsNotExist(err)
}

/**
* 根据文件名称创建一个 os.File 对象
* 注意；会取 filename我前2位作为目录名称:方便后续查找
*
* */
func createFile(hash string) (*os.File, error) {
	dir, filepath := MakeOssStorageFilePath(hash)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}

	return os.Create(filepath)
}
