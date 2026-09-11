package response

import (
	"fmt"
	"net/http"
	"reflect"
)

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

func setFiled(object reflect.Value, key string, value string) {
	filed := object.FieldByName(key)

	if !filed.IsValid() {
		panic(fmt.Sprintf("field '%s' is not valid", key))
	}

	if !filed.CanSet() {
		panic(fmt.Sprintf("field '%s' cannot be set", key))
	}

	filed.Set(reflect.ValueOf(value))
}

func MakeResponseWithUserInfo[T any](data T) *KResponse {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		panic("data must be a pointer")
		return nil
	}

	v = v.Elem()

	setFiled(v, "NickName", "王恒星")
	setFiled(v, "Avatar", "王恒星")

	return &KResponse{
		Data: data,
		Code: http.StatusOK,
	}
}
