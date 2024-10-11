package _struct

// LoginRequest 登录请求结构体
type LoginRequest struct {
	UserName string
	Password string
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	UserName string
	UserID   int
	Token    string
}

// RegisterRequest 用户注册结构体
type RegisterRequest struct {
	UserName string
	Password string
}

// RegisterResponse 注册响应结构体
type RegisterResponse struct {
	UserID   int
	UserName string
}
