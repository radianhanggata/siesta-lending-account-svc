package internalerror

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *Response) Error() string {
	return r.Message
}

var ErrNotFound = &Response{Code: 404, Message: "not found"}
var ErrBadRequest = &Response{Code: 400, Message: "bad request"}
var ErrInternalServer = &Response{Code: 500, Message: "internal server error"}
