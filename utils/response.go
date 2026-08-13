package utils

import (
	"net/http"
)

type KResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code"`
}

func MakeResponse(data any) *KResponse {
	return &KResponse{
		Data: data,
		Code: http.StatusOK,
	}
}

func MakeResponseWidthCode(data any, code int) *KResponse {
	return &KResponse{
		Data: data,
		Code: code,
	}
}
