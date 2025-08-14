package dto

const (
	statusOK          = "OK"
	statusError       = "ERROR"
	internalErrorText = "Internal Error"
)

type Response struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func OK() *Response {
	return &Response{
		Status: statusOK,
	}
}

func Error(err error) *Response {
	return &Response{
		Status: statusError,
		Error:  err.Error(),
	}
}

func InternalError() *Response {
	return &Response{
		Status: statusError,
		Error:  internalErrorText,
	}
}
