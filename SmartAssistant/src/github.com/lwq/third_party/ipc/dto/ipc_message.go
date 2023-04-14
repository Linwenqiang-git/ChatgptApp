package dto

import (
	"github.com/google/uuid"
	. "github.com/lwq/internal/shared/consts"
)

type IpcRequest struct {
	Id      uuid.UUID
	Module  AppModule
	Message string
	IsExit  bool
}

func ExitIpcRequest() IpcRequest {
	return IpcRequest{
		IsExit: true,
	}
}
func OpenaiKeyIpcRequest(openaiKey string) IpcRequest {
	return IpcRequest{
		Id:      uuid.New(),
		Module:  OpenaiKey,
		Message: openaiKey, //直接将api key发送，python端接收后无法解析
	}
}

type IpcResponse struct {
	ResponseId uuid.UUID
	Code       int
	Message    string
	Error      error
	ErrorMsg   string
}

func IpcResponseError(err error, msg string) IpcResponse {
	return IpcResponse{
		Code:     500,
		Error:    err,
		ErrorMsg: msg,
	}
}
