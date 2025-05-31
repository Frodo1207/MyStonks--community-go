package response

type Response struct {
	Code ErrorCode   `json:"code"`
	Args []string    `json:"args"`
	Data interface{} `json:"data"`
}

type ErrorCode int

const (
	// general
	ErrorCodeSuccess          ErrorCode = 0
	ErrorCodeInternalError    ErrorCode = 10001
	ErrorCodeInvalidToken     ErrorCode = 10002
	ErrorCodeInvalidRequest   ErrorCode = 10003
	ErrorCodeInvalidSignature ErrorCode = 10004

	// user
	ErrorCodeUserNotFound ErrorCode = 20001
)

func SuccessResponse(data interface{}) Response {
	return Response{
		Code: ErrorCodeSuccess,
		Args: []string{},
		Data: data,
	}
}

func ErrorResponse(code ErrorCode, args []string) Response {
	return Response{
		Code: code,
		Args: args,
		Data: map[string]interface{}{},
	}
}
