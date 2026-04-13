package response

import (
	"context"
	"net/http"
	"strings"

	"github.com/FatSong207/common/go-zero/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 統一回應結構
type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func init() {
	// 全局錯誤攔截器
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		var body Body
		var statusCode int

		// 記錄詳細錯誤到 Log 檔案與終端機，免得被吃掉
		logx.Errorf("API Request Error: %+v", err)

		// 先嘗試看是不是自訂業務錯誤 CodeError
		switch e := err.(type) {
		case *errorx.CodeError:
			body.Code = e.Code
			body.Msg = e.Msg
			statusCode = e.HttpStatusCode()
		default:
			// 再嘗試看是不是從 RPC 傳過來的 gRPC 錯誤
			if rpcErr := errorx.FromGrpcError(err); rpcErr != nil && rpcErr.Code != errorx.ServerError {
				body.Code = rpcErr.Code
				body.Msg = rpcErr.Msg
				statusCode = rpcErr.HttpStatusCode()
			} else {
				// 處理 HTTP Request Body 缺失、格式解析失敗、或欄位驗證不過 (這不應該被當成內部 500 錯誤)
				errStr := err.Error()
				lowerErrStr := strings.ToLower(errStr)
				isParamError := errStr == "EOF" ||
					strings.Contains(lowerErrStr, "json") ||
					strings.Contains(lowerErrStr, "parse") ||
					strings.Contains(lowerErrStr, "is not set") ||
					strings.Contains(lowerErrStr, "invalid") ||
					strings.Contains(lowerErrStr, "mismatch")

				if isParamError {
					body.Code = errorx.ParamError
					body.Msg = "Invalid Param： " + errStr
					statusCode = http.StatusBadRequest
				} else {
					// 真正的未知錯誤
					body.Code = errorx.ServerError
					body.Msg = "Server Error"
					statusCode = http.StatusInternalServerError
				}
			}
		}
		// 返回對應的 http status code, 與 json body
		return statusCode, body
	})

	// 全局成功攔截器
	httpx.SetOkHandler(func(ctx context.Context, data interface{}) interface{} {
		return Body{
			Code: errorx.OK,
			Msg:  "success",
			Data: data,
		}
	})
}
