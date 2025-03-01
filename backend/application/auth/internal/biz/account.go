package biz

// 集成登录
type (
	SigninRequest struct {
		Code  string `json:"code,omitempty"`
		State string `json:"state,omitempty"`
	}
	SigninReply struct {
		State string `json:"state,omitempty"`
		Data  string `json:"data,omitempty"`
	}
)
