package libs

// HttpResponse 封装接口返回类型
type HttpResponse struct {
	Code     int         `json:"code"`
	HttpData interface{} `json:"data"`
	Message  string      `json:"message"`
}

// 公共方法
func response(code int, data interface{}, message string) *HttpResponse {
	h := &HttpResponse{}
	h.Code = code
	h.HttpData = data
	h.Message = message
	return h
}

// Success 请求成功，返回状态码为200， 消息， 数据
func Success(data interface{}, message string) *HttpResponse {
	return response(200, data, message)
}

// ServerError 服务器异常，返回状态码为500， 消息
func ServerError(message string) *HttpResponse {
	return response(500, nil, message)
}

// ParamsError 参数异常， 返回状态码400， 消息
func ParamsError(message string) *HttpResponse {
	return response(400, nil, message)
}

// SessionError 权限异常，返回状态码
func SessionError(message string) *HttpResponse {
	return response(401, nil, message)
}

//
func AuthorError(message string) *HttpResponse {
	return response(405, nil, message)
}
