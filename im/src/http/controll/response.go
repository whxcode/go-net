package controll

import "go-net/im/src/http/consts"

type KResponse struct {
	Data any
	Code consts.Status
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
