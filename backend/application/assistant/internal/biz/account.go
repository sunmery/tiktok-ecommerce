package biz

import "github.com/casdoor/casdoor-go-sdk/casdoorsdk"

type SigninRequest struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}
type SigninReply struct {
	State string `json:"state,omitempty"`
	Data  string `json:"data,omitempty"`
}
type GetUserInfoRequest struct {
	Authorization string
}

type GetUserInfoReply struct {
	State string          `json:"state,omitempty"`
	Data  casdoorsdk.User `json:"data"`
}
