package response

import (
	"context"
	"net/http"
	"strings"

	"github.com/FatSong207/common/go-zero/errorx"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body unified API response structure
type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func init() {
	// Global error handler
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		var body Body
		var statusCode int

		// Log detailed error for debugging purposes
		logx.Errorf("API request error: %+v", err)

		// Try custom business error first
		switch e := err.(type) {
		case *errorx.CodeError:
			body.Code = e.Code
			body.Msg = e.Msg
			statusCode = e.HttpStatusCode()
		default:
			// Check if it's a gRPC error from RPC calls
			if rpcErr := errorx.FromGrpcError(err); rpcErr != nil && rpcErr.Code != errorx.ServerError {
				body.Code = rpcErr.Code
				body.Msg = rpcErr.Msg
				statusCode = rpcErr.HttpStatusCode()
			} else {
				// Handle request body, parsing, and validation errors (not 500)
				if isParamError(err) {
					body.Code = errorx.ParamError
					body.Msg = "Invalid request parameters"
					statusCode = http.StatusBadRequest
				} else {
					// Unknown internal error
					body.Code = errorx.ServerError
					body.Msg = "Internal server error"
					statusCode = http.StatusInternalServerError
				}
			}
		}
		return statusCode, body
	})

	// Global success handler
	httpx.SetOkHandler(func(ctx context.Context, data interface{}) interface{} {
		return Body{
			Code: errorx.OK,
			Msg:  "Success",
			Data: data,
		}
	})
}

// isParamError checks if the error is a parameter validation or parsing error
func isParamError(err error) bool {
	errStr := err.Error()
	lowerErrStr := strings.ToLower(errStr)

	paramErrorIndicators := []string{
		"EOF",
		"json",
		"parse",
		"is not set",
		"invalid",
		"mismatch",
		"required",
		"decode",
		"unmarshal",
	}

	if errStr == "EOF" {
		return true
	}

	for _, indicator := range paramErrorIndicators {
		if strings.Contains(lowerErrStr, indicator) {
			return true
		}
	}

	return false
}
