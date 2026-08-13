package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Level uint32

const (
	Log   Level = 1 << 0
	Info        = 1 << 1
	WARN        = 1 << 2
	ERROR       = 1 << 3
)

type Config struct {
	Server struct {
		Port         string `json:"port"`         // 端口
		TemplatePath string `json:"templatePath"` // 前端模版路径
		FileOss      string `json:"fileOss"`      // 文件oss路径
	} `json:"server"`

	Log struct {
		PathDir string `json:"pathDir"` // 日志目录
		Level   Level  `json:"level"`   // 日志级别
	} `json:"log"`
}

func (c *Config) Dump() {
	fmt.Println("============ config ==================")
	fmt.Printf("Server Port: %s\n", c.Server.Port)
	fmt.Printf("Server Template Path: %s\n", c.Server.TemplatePath)
	fmt.Printf("Server FileOss Path: %s\n", c.Server.FileOss)

	fmt.Printf("Log PathDir: %s\n", c.Log.PathDir)
	fmt.Printf("Log Level: %d\n", c.Log.Level)
	fmt.Println("============ config ==================")
}

var ConfigData *Config

func init() {
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	ConfigData = &Config{}

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(ConfigData); err != nil {
		panic(err)
	}
}
