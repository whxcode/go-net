package controll

import "go-net/http/consts"

type KResponse struct {
	Data any           `json:"data"`
	Code consts.Status `json:"code"`
}

func MakeResponse(data any) *KResponse {
	return &KResponse{
		Data: data,
		Code: consts.KSuccess,
	}
}

func MakeResponseWidthCode(data any, code consts.Status) *KResponse {
	return &KResponse{
		Data: data,
		Code: code,
	}
}
