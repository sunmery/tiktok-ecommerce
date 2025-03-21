// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.3
// - protoc             v5.29.3
// source: v1/user.proto

package userv1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationUserServiceCreateAddresses = "/ecommerce.user.v1.UserService/CreateAddresses"
const OperationUserServiceCreateCreditCard = "/ecommerce.user.v1.UserService/CreateCreditCard"
const OperationUserServiceDeleteAddresses = "/ecommerce.user.v1.UserService/DeleteAddresses"
const OperationUserServiceDeleteCreditCard = "/ecommerce.user.v1.UserService/DeleteCreditCard"
const OperationUserServiceGetAddress = "/ecommerce.user.v1.UserService/GetAddress"
const OperationUserServiceGetAddresses = "/ecommerce.user.v1.UserService/GetAddresses"
const OperationUserServiceGetCreditCard = "/ecommerce.user.v1.UserService/GetCreditCard"
const OperationUserServiceGetUserProfile = "/ecommerce.user.v1.UserService/GetUserProfile"
const OperationUserServiceListCreditCards = "/ecommerce.user.v1.UserService/ListCreditCards"
const OperationUserServiceUpdateAddresses = "/ecommerce.user.v1.UserService/UpdateAddresses"

type UserServiceHTTPServer interface {
	// CreateAddresses 创建用户地址
	CreateAddresses(context.Context, *Address) (*Address, error)
	// CreateCreditCard 创建用户的信用卡信息
	CreateCreditCard(context.Context, *CreditCard) (*emptypb.Empty, error)
	// DeleteAddresses 删除用户地址
	DeleteAddresses(context.Context, *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	// DeleteCreditCard 删除用户的信用卡信息
	DeleteCreditCard(context.Context, *DeleteCreditCardsRequest) (*emptypb.Empty, error)
	// GetAddress 根据 ID获取用户地址
	GetAddress(context.Context, *GetAddressRequest) (*Address, error)
	// GetAddresses 获取用户地址列表
	GetAddresses(context.Context, *emptypb.Empty) (*GetAddressesReply, error)
	// GetCreditCard 根据ID搜索用户的信用卡信息
	GetCreditCard(context.Context, *GetCreditCardRequest) (*CreditCard, error)
	// GetUserProfile 获取用户个人资料
	GetUserProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
	// ListCreditCards 列出用户的信用卡信息
	ListCreditCards(context.Context, *emptypb.Empty) (*CreditCards, error)
	// UpdateAddresses 更新用户地址
	UpdateAddresses(context.Context, *Address) (*Address, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/users/profile", _UserService_GetUserProfile0_HTTP_Handler(srv))
	r.POST("/v1/users/address", _UserService_CreateAddresses0_HTTP_Handler(srv))
	r.PATCH("/v1/users/address", _UserService_UpdateAddresses0_HTTP_Handler(srv))
	r.DELETE("/v1/users/address", _UserService_DeleteAddresses0_HTTP_Handler(srv))
	r.GET("/v1/users/address", _UserService_GetAddress0_HTTP_Handler(srv))
	r.POST("/v1/users/credit_cards", _UserService_CreateCreditCard0_HTTP_Handler(srv))
	r.GET("/v1/users/addresses", _UserService_GetAddresses0_HTTP_Handler(srv))
	r.GET("/v1/users/credit_cards", _UserService_ListCreditCards0_HTTP_Handler(srv))
	r.GET("/v1/users/credit_cards/{id}", _UserService_GetCreditCard0_HTTP_Handler(srv))
	r.DELETE("/v1/users/credit_cards/{id}", _UserService_DeleteCreditCard0_HTTP_Handler(srv))
}

func _UserService_GetUserProfile0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetProfileRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetUserProfile)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserProfile(ctx, req.(*GetProfileRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetProfileResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_CreateAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in Address
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceCreateAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateAddresses(ctx, req.(*Address))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Address)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in Address
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateAddresses(ctx, req.(*Address))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Address)
		return ctx.Result(200, reply)
	}
}

func _UserService_DeleteAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteAddressesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceDeleteAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteAddresses(ctx, req.(*DeleteAddressesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteAddressesReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_GetAddress0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetAddressRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetAddress)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAddress(ctx, req.(*GetAddressRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Address)
		return ctx.Result(200, reply)
	}
}

func _UserService_CreateCreditCard0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreditCard
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceCreateCreditCard)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateCreditCard(ctx, req.(*CreditCard))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _UserService_GetAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAddresses(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAddressesReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_ListCreditCards0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceListCreditCards)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCreditCards(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreditCards)
		return ctx.Result(200, reply)
	}
}

func _UserService_GetCreditCard0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCreditCardRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetCreditCard)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCreditCard(ctx, req.(*GetCreditCardRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreditCard)
		return ctx.Result(200, reply)
	}
}

func _UserService_DeleteCreditCard0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteCreditCardsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceDeleteCreditCard)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteCreditCard(ctx, req.(*DeleteCreditCardsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type UserServiceHTTPClient interface {
	CreateAddresses(ctx context.Context, req *Address, opts ...http.CallOption) (rsp *Address, err error)
	CreateCreditCard(ctx context.Context, req *CreditCard, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteAddresses(ctx context.Context, req *DeleteAddressesRequest, opts ...http.CallOption) (rsp *DeleteAddressesReply, err error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetAddress(ctx context.Context, req *GetAddressRequest, opts ...http.CallOption) (rsp *Address, err error)
	GetAddresses(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *GetAddressesReply, err error)
	GetCreditCard(ctx context.Context, req *GetCreditCardRequest, opts ...http.CallOption) (rsp *CreditCard, err error)
	GetUserProfile(ctx context.Context, req *GetProfileRequest, opts ...http.CallOption) (rsp *GetProfileResponse, err error)
	ListCreditCards(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *CreditCards, err error)
	UpdateAddresses(ctx context.Context, req *Address, opts ...http.CallOption) (rsp *Address, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) CreateAddresses(ctx context.Context, in *Address, opts ...http.CallOption) (*Address, error) {
	var out Address
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceCreateAddresses))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) CreateCreditCard(ctx context.Context, in *CreditCard, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/users/credit_cards"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceCreateCreditCard))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) DeleteAddresses(ctx context.Context, in *DeleteAddressesRequest, opts ...http.CallOption) (*DeleteAddressesReply, error) {
	var out DeleteAddressesReply
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceDeleteAddresses))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) DeleteCreditCard(ctx context.Context, in *DeleteCreditCardsRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/users/credit_cards/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceDeleteCreditCard))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetAddress(ctx context.Context, in *GetAddressRequest, opts ...http.CallOption) (*Address, error) {
	var out Address
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetAddress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetAddresses(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*GetAddressesReply, error) {
	var out GetAddressesReply
	pattern := "/v1/users/addresses"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetAddresses))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetCreditCard(ctx context.Context, in *GetCreditCardRequest, opts ...http.CallOption) (*CreditCard, error) {
	var out CreditCard
	pattern := "/v1/users/credit_cards/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetCreditCard))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetUserProfile(ctx context.Context, in *GetProfileRequest, opts ...http.CallOption) (*GetProfileResponse, error) {
	var out GetProfileResponse
	pattern := "/v1/users/profile"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetUserProfile))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) ListCreditCards(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*CreditCards, error) {
	var out CreditCards
	pattern := "/v1/users/credit_cards"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceListCreditCards))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateAddresses(ctx context.Context, in *Address, opts ...http.CallOption) (*Address, error) {
	var out Address
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateAddresses))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PATCH", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
