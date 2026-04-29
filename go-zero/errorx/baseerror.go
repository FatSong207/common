package errorx

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CodeError unified error structure
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// HttpStatusCode maps error code to corresponding HTTP status code
func (e *CodeError) HttpStatusCode() int {
	switch e.Code {
	case ParamError:
		return http.StatusBadRequest // 400
	case AuthError:
		return http.StatusUnauthorized // 401
	case PermissionDenied:
		return http.StatusForbidden // 403
	case ResourceNotFound:
		return http.StatusNotFound // 404
	case Conflict:
		return http.StatusConflict // 409
	case TooManyRequests:
		return http.StatusTooManyRequests // 429
	case ServerError:
		return http.StatusInternalServerError // 500
	case ServiceUnavailable:
		return http.StatusServiceUnavailable // 503
	default:
		return http.StatusInternalServerError // Default 500
	}
}

// Err returns a gRPC standard error for RPC transmission
func (e *CodeError) Err() error {
	return status.Error(codes.Code(e.Code), e.Msg)
}

// FromGrpcError converts gRPC error back to CodeError for API layer
func FromGrpcError(err error) *CodeError {
	if err == nil {
		return nil
	}
	// Parse gRPC status
	if st, ok := status.FromError(err); ok {
		return &CodeError{
			Code: int(st.Code()),
			Msg:  st.Message(),
		}
	}
	// Return server error if unable to parse
	return NewServerError()
}

// Error code definitions
const (
	// OK = 0 成功 HTTP 200
	// 用途: 业务正常完成时返回
	// 示例: return nil (框架自动返回 OK)
	OK = 0

	// ServerError = 1001 服务器内部错误 HTTP 500
	// 用途: 未预期的系统错误、数据库连接失败、依赖服务故障等
	// 便利: NewServerError() 直接返回
	ServerError = 1001

	// ParamError = 1002 请求参数错误 HTTP 400
	// 用途: JSON解析失败、字段格式错误、必填字段缺失、字段类型不匹配等
	// 便利: NewParamError(msg) 可自定义具体错误内容
	ParamError = 1002

	// AuthError = 1003 认证失败 HTTP 401
	// 用途: 用户未登录、Token过期、Token无效、登录信息丢失等
	// 便利: NewAuthError() 直接返回
	AuthError = 1003

	// PermissionDenied = 1004 权限不足 HTTP 403
	// 用途: 用户已登录但无权限访问该资源、无权限执行该操作等
	// 便利: NewPermissionDenied(msg) 可指定权限要求
	PermissionDenied = 1004

	// ResourceNotFound = 1005 资源不存在 HTTP 404
	// 用途: 查询的用户/订单/文件等不存在、记录已删除、URL错误等
	// 便利: NewResourceNotFound(resource) 自动拼接 "xxx not found"
	ResourceNotFound = 1005

	// Conflict = 1006 资源冲突 HTTP 409
	// 用途: 唯一约束冲突（重复注册邮箱）、版本冲突、状态冲突等
	// 便利: NewConflict(msg) 可描述冲突详情
	Conflict = 1006

	// TooManyRequests = 1007 请求过于频繁 HTTP 429
	// 用途: API速率限制、操作频率限制（如1分钟内只能发送3条短信）等
	// 便利: NewTooManyRequests() 直接返回
	TooManyRequests = 1007

	// ValidationError = 1008 数据验证失败 HTTP 500
	// 用途: 业务逻辑验证失败，如余额不足、库存不足、年龄不符合要求等
	// 便利: NewValidationError(msg) 可指定验证失败原因
	ValidationError = 1008

	// BusinessError = 1009 业务逻辑错误 HTTP 500
	// 用途: 其他业务相关错误，如转账失败、操作不合法、状态错误等
	// 便利: NewBusinessError(msg) 通用业务错误处理
	BusinessError = 1009

	// ServiceUnavailable = 1010 服务不可用 HTTP 503
	// 用途: 依赖服务故障、数据库维护中、服务升级中等临时不可用
	// 便利: NewServiceUnavailable() 直接返回
	ServiceUnavailable = 1010
)

var codeMsg = map[int]string{
	OK:                 "Success",
	ServerError:        "Internal server error",
	ParamError:         "Invalid request parameters",
	AuthError:          "Unauthorized",
	PermissionDenied:   "Permission denied",
	ResourceNotFound:   "Resource not found",
	Conflict:           "Resource conflict",
	TooManyRequests:    "Too many requests",
	ValidationError:    "Validation failed",
	BusinessError:      "Business logic error",
	ServiceUnavailable: "Service temporarily unavailable",
}

// NewCodeError 根据错误码创建错误，使用预定义的消息
// 用途: 当错误码对应的默认消息已满足需求时使用
// 示例: return errorx.NewCodeError(errorx.AuthError)
func NewCodeError(code int) *CodeError {
	return &CodeError{Code: code, Msg: codeMsg[code]}
}

// NewCodeErrorMsg 根据错误码和自定义消息创建错误
// 用途: 需要自定义错误消息时使用
// 示例: return errorx.NewCodeErrorMsg(errorx.ParamError, "name is required")
func NewCodeErrorMsg(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

// NewServerError 创建服务器内部错误
// 用途: 捕获未处理的异常或系统错误时返回
// 示例: return errorx.NewServerError()
func NewServerError() *CodeError {
	return NewCodeError(ServerError)
}

// NewParamError 创建参数验证错误
// 用途: 请求参数不合法时返回，可自定义具体错误信息
// 示例: return errorx.NewParamError("email format invalid")
func NewParamError(msg string) *CodeError {
	return NewCodeErrorMsg(ParamError, msg)
}

// NewAuthError 创建认证失败错误
// 用途: 用户登录状态失效或未登录时返回
// 示例: return errorx.NewAuthError()
func NewAuthError() *CodeError {
	return NewCodeError(AuthError)
}

// NewPermissionDenied 创建权限不足错误
// 用途: 用户无权限访问资源或执行操作时返回
// 示例: return errorx.NewPermissionDenied("only admin can delete users")
func NewPermissionDenied(msg string) *CodeError {
	return NewCodeErrorMsg(PermissionDenied, msg)
}

// NewResourceNotFound 创建资源不存在错误
// 用途: 查询的资源不存在或已删除时返回，自动拼接 "not found"
// 示例: return errorx.NewResourceNotFound("user") // 返回 "user not found"
func NewResourceNotFound(resource string) *CodeError {
	return NewCodeErrorMsg(ResourceNotFound, resource+" not found")
}

// NewConflict 创建资源冲突错误
// 用途: 资源重复、状态冲突、唯一约束冲突等
// 示例: return errorx.NewConflict("email already registered")
func NewConflict(msg string) *CodeError {
	return NewCodeErrorMsg(Conflict, msg)
}

// NewTooManyRequests 创建请求频率过高错误
// 用途: 触发API限流或操作频率限制时返回
// 示例: return errorx.NewTooManyRequests()
func NewTooManyRequests() *CodeError {
	return NewCodeError(TooManyRequests)
}

// NewValidationError 创建数据验证失败错误
// 用途: 业务逻辑验证失败（非参数格式错误）时返回
// 示例: return errorx.NewValidationError("insufficient balance")
func NewValidationError(msg string) *CodeError {
	return NewCodeErrorMsg(ValidationError, msg)
}

// NewBusinessError 创建业务逻辑错误
// 用途: 其他业务相关错误时返回，是最通用的业务错误
// 示例: return errorx.NewBusinessError("failed to process payment")
func NewBusinessError(msg string) *CodeError {
	return NewCodeErrorMsg(BusinessError, msg)
}

// NewServiceUnavailable 创建服务不可用错误
// 用途: 依赖服务故障或服务维护中时返回
// 示例: return errorx.NewServiceUnavailable()
func NewServiceUnavailable() *CodeError {
	return NewCodeError(ServiceUnavailable)
}
