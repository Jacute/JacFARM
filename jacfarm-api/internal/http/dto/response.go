package dto

const (
	statusOK    = "OK"
	statusError = "ERROR"
)

var (
	ErrDecodingBody       = Error("error decoding body")
	ErrInvalidContentType = Error("invalid content-type")
	ErrInternal           = Error("internal error")
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

func Error(err string) *Response {
	return &Response{
		Status: statusError,
		Error:  err,
	}
}
