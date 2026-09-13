package response

import (
	"net/http"
	"reflect"

	"go-net/model"
)

type IUserPair interface {
	SetNickName(nickname string)
	SetAvatar(avatar string)
}

type KResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code"`
}

func MakeResponse[T any](data T) *KResponse {
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

func MakeResponseWithUserInfo[T IUserPair](data T) *KResponse {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName("UserPair")

	/*
		fmt.Println("field.Type:", field.Type())
		fmt.Println("field.Kind:", field.Kind())
		fmt.Println("field.IsNil:", field.IsNil())
	*/

	field.Set(reflect.ValueOf(&model.UserPair{
		Nickname: "whx",
		Avatar:   "https://example.com/avatar.png",
	}))

	return &KResponse{
		Data: data,
		Code: http.StatusOK,
	}
}
