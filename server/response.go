// Package utils
// @Description 封装返回的数据格式
package server

// ResultJson 接口返回数据的结构体
type ResultJson struct {
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

// RequestSuccess 请求成功时的响应
func RequestSuccess(result interface{}) ResultJson {
	return ResultJson{Code: 200, Type: "success", Result: result, Message: "请求成功"}
}

// RequestFail 请求失败时的响应
func RequestFail(error string) ResultJson {
	return ResultJson{Code: 500, Type: "error", Result: nil, Message: error}
}

// RequestNoPermission 没有权限时的响应
func RequestNoPermission() ResultJson {
	return ResultJson{Code: 401, Type: "warning", Result: nil, Message: "无权访问"}
}
