package errorx

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CodeError 統一錯誤結構
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// HttpStatusCode 根據錯誤碼返回對應的 HTTP 狀態碼
func (e *CodeError) HttpStatusCode() int {
	switch e.Code {
	case ParamError:
		return http.StatusBadRequest // 400
	case AuthError:
		return http.StatusUnauthorized // 401
	case UserNotFound:
		return http.StatusNotFound // 404
	case ServerError:
		return http.StatusInternalServerError // 500
	default:
		return http.StatusInternalServerError // 預設 500
	}
}

// Err 回傳 gRPC 標準錯誤 (用來在 RPC 層回傳給呼叫端)
// 這樣做可以確保傳輸時 statusCode 與 錯誤訊息 不會遺失
func (e *CodeError) Err() error {
	return status.Error(codes.Code(e.Code), e.Msg)
}

// FromGrpcError 把 gRPC 取回的錯誤還原成 CodeError (給 API 層使用)
func FromGrpcError(err error) *CodeError {
	if err == nil {
		return nil
	}
	// 解析 gRPC 的 status 特性
	if st, ok := status.FromError(err); ok {
		return &CodeError{
			Code: int(st.Code()),
			Msg:  st.Message(),
		}
	}
	// 如果無法解析或未知錯誤，當作普通的系統錯誤
	return NewServerError()
}

// 錯誤碼定義
const (
	OK           = 0
	ServerError  = 1001
	ParamError   = 1002
	AuthError    = 1003
	UserNotFound = 1004
)

var codeMsg = map[int]string{
	OK:           "success",
	ServerError:  "server error",
	ParamError:   "invalid param",
	AuthError:    "Unauthorized",
	UserNotFound: "user not found",
}

// NewCodeError 建立錯誤
func NewCodeError(code int) *CodeError {
	return &CodeError{Code: code, Msg: codeMsg[code]}
}

// NewCodeErrorMsg 建立自訂訊息錯誤
func NewCodeErrorMsg(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

// NewServerError 建立伺服器錯誤
func NewServerError() *CodeError {
	return NewCodeError(ServerError)
}

// NewParamError 建立參數錯誤
func NewParamError(msg string) *CodeError {
	return NewCodeErrorMsg(ParamError, msg)
}
