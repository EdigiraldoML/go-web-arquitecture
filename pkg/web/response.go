package web

import "strconv"

type Response struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func NewResponse(code int, data interface{}, err string) (r Response) {
	r.Code = strconv.FormatInt(int64(code), 10)

	if code < 300 {
		r.Data = data
		r.Error = ""

		return r
	}
	r.Data = nil
	r.Error = err

	return r
}
