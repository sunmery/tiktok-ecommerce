package biz

import "github.com/casdoor/casdoor-go-sdk/casdoorsdk"

type GetProfileRequest struct {
	Authorization string
}

type GetProfileReply struct {
	State string          `json:"state,omitempty"`
	Data  casdoorsdk.User `json:"data"`
}
