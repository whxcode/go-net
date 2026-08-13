package utils

import "time"

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetToday() string {
	return time.Now().Format("2006-01-02")
}
