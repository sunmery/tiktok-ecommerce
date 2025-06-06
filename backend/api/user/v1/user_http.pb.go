// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
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

const OperationUserServiceCreateConsumerAddresses = "/ecommerce.user.v1.UserService/CreateConsumerAddresses"
const OperationUserServiceCreateCreditCard = "/ecommerce.user.v1.UserService/CreateCreditCard"
const OperationUserServiceDeleteConsumerAddresses = "/ecommerce.user.v1.UserService/DeleteConsumerAddresses"
const OperationUserServiceDeleteCreditCard = "/ecommerce.user.v1.UserService/DeleteCreditCard"
const OperationUserServiceDeleteFavorites = "/ecommerce.user.v1.UserService/DeleteFavorites"
const OperationUserServiceDeleteUser = "/ecommerce.user.v1.UserService/DeleteUser"
const OperationUserServiceGetConsumerAddress = "/ecommerce.user.v1.UserService/GetConsumerAddress"
const OperationUserServiceGetConsumerAddresses = "/ecommerce.user.v1.UserService/GetConsumerAddresses"
const OperationUserServiceGetCreditCard = "/ecommerce.user.v1.UserService/GetCreditCard"
const OperationUserServiceGetFavorites = "/ecommerce.user.v1.UserService/GetFavorites"
const OperationUserServiceGetUserProfile = "/ecommerce.user.v1.UserService/GetUserProfile"
const OperationUserServiceGetUsers = "/ecommerce.user.v1.UserService/GetUsers"
const OperationUserServiceListCreditCards = "/ecommerce.user.v1.UserService/ListCreditCards"
const OperationUserServiceSetFavorites = "/ecommerce.user.v1.UserService/SetFavorites"
const OperationUserServiceUpdateConsumerAddresses = "/ecommerce.user.v1.UserService/UpdateConsumerAddresses"
const OperationUserServiceUpdateUser = "/ecommerce.user.v1.UserService/UpdateUser"

type UserServiceHTTPServer interface {
	// CreateConsumerAddresses 创建用户地址
	CreateConsumerAddresses(context.Context, *ConsumerAddress) (*ConsumerAddress, error)
	// CreateCreditCard 创建用户的信用卡信息
	CreateCreditCard(context.Context, *CreditCard) (*emptypb.Empty, error)
	// DeleteConsumerAddresses 删除用户地址
	DeleteConsumerAddresses(context.Context, *DeleteConsumerAddressesRequest) (*DeleteConsumerAddressesReply, error)
	// DeleteCreditCard 删除用户的信用卡信息
	DeleteCreditCard(context.Context, *DeleteCreditCardsRequest) (*emptypb.Empty, error)
	// DeleteFavorites 删除商品收藏
	DeleteFavorites(context.Context, *UpdateFavoritesRequest) (*UpdateFavoritesResply, error)
	// DeleteUser 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	// GetConsumerAddress 根据 ID获取用户地址
	GetConsumerAddress(context.Context, *GetConsumerAddressRequest) (*ConsumerAddress, error)
	// GetConsumerAddresses 获取用户地址列表
	GetConsumerAddresses(context.Context, *emptypb.Empty) (*GetConsumerAddressesReply, error)
	// GetCreditCard 根据ID搜索用户的信用卡信息
	GetCreditCard(context.Context, *GetCreditCardRequest) (*CreditCard, error)
	// GetFavorites 获取用户商品收藏
	GetFavorites(context.Context, *GetFavoritesRequest) (*Favorites, error)
	// GetUserProfile 获取用户个人资料
	GetUserProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
	// GetUsers 获取全部用户信息
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	// ListCreditCards 列出用户的信用卡信息
	ListCreditCards(context.Context, *emptypb.Empty) (*CreditCards, error)
	// SetFavorites 添加商品收藏
	SetFavorites(context.Context, *UpdateFavoritesRequest) (*UpdateFavoritesResply, error)
	// UpdateConsumerAddresses 更新用户地址
	UpdateConsumerAddresses(context.Context, *ConsumerAddress) (*ConsumerAddress, error)
	// UpdateUser 更新用户信息
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/users/profile", _UserService_GetUserProfile0_HTTP_Handler(srv))
	r.GET("/v1/users", _UserService_GetUsers0_HTTP_Handler(srv))
	r.POST("/v1/users", _UserService_DeleteUser0_HTTP_Handler(srv))
	r.POST("/v1/users/address", _UserService_CreateConsumerAddresses0_HTTP_Handler(srv))
	r.PATCH("/v1/users/address", _UserService_UpdateConsumerAddresses0_HTTP_Handler(srv))
	r.DELETE("/v1/users/address", _UserService_DeleteConsumerAddresses0_HTTP_Handler(srv))
	r.GET("/v1/users/address", _UserService_GetConsumerAddress0_HTTP_Handler(srv))
	r.POST("/v1/users/credit_cards", _UserService_CreateCreditCard0_HTTP_Handler(srv))
	r.GET("/v1/users/addresses", _UserService_GetConsumerAddresses0_HTTP_Handler(srv))
	r.GET("/v1/users/credit_cards", _UserService_ListCreditCards0_HTTP_Handler(srv))
	r.GET("/v1/users/favorites", _UserService_GetFavorites0_HTTP_Handler(srv))
	r.PUT("/v1/users/favorites", _UserService_SetFavorites0_HTTP_Handler(srv))
	r.DELETE("/v1/users/favorites", _UserService_DeleteFavorites0_HTTP_Handler(srv))
	r.POST("/v1/users/{user_id}", _UserService_UpdateUser0_HTTP_Handler(srv))
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

func _UserService_GetUsers0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetUsersRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetUsers)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUsers(ctx, req.(*GetUsersRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetUsersResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_DeleteUser0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceDeleteUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteUser(ctx, req.(*DeleteUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteUserResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_CreateConsumerAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ConsumerAddress
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceCreateConsumerAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateConsumerAddresses(ctx, req.(*ConsumerAddress))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ConsumerAddress)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateConsumerAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ConsumerAddress
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateConsumerAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateConsumerAddresses(ctx, req.(*ConsumerAddress))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ConsumerAddress)
		return ctx.Result(200, reply)
	}
}

func _UserService_DeleteConsumerAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteConsumerAddressesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceDeleteConsumerAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteConsumerAddresses(ctx, req.(*DeleteConsumerAddressesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteConsumerAddressesReply)
		return ctx.Result(200, reply)
	}
}

func _UserService_GetConsumerAddress0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetConsumerAddressRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetConsumerAddress)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetConsumerAddress(ctx, req.(*GetConsumerAddressRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ConsumerAddress)
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

func _UserService_GetConsumerAddresses0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetConsumerAddresses)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetConsumerAddresses(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetConsumerAddressesReply)
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

func _UserService_GetFavorites0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetFavoritesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceGetFavorites)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetFavorites(ctx, req.(*GetFavoritesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Favorites)
		return ctx.Result(200, reply)
	}
}

func _UserService_SetFavorites0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateFavoritesRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceSetFavorites)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SetFavorites(ctx, req.(*UpdateFavoritesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateFavoritesResply)
		return ctx.Result(200, reply)
	}
}

func _UserService_DeleteFavorites0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateFavoritesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceDeleteFavorites)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteFavorites(ctx, req.(*UpdateFavoritesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateFavoritesResply)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateUser0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUser(ctx, req.(*UpdateUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateUserResponse)
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
	CreateConsumerAddresses(ctx context.Context, req *ConsumerAddress, opts ...http.CallOption) (rsp *ConsumerAddress, err error)
	CreateCreditCard(ctx context.Context, req *CreditCard, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteConsumerAddresses(ctx context.Context, req *DeleteConsumerAddressesRequest, opts ...http.CallOption) (rsp *DeleteConsumerAddressesReply, err error)
	DeleteCreditCard(ctx context.Context, req *DeleteCreditCardsRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	DeleteFavorites(ctx context.Context, req *UpdateFavoritesRequest, opts ...http.CallOption) (rsp *UpdateFavoritesResply, err error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest, opts ...http.CallOption) (rsp *DeleteUserResponse, err error)
	GetConsumerAddress(ctx context.Context, req *GetConsumerAddressRequest, opts ...http.CallOption) (rsp *ConsumerAddress, err error)
	GetConsumerAddresses(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *GetConsumerAddressesReply, err error)
	GetCreditCard(ctx context.Context, req *GetCreditCardRequest, opts ...http.CallOption) (rsp *CreditCard, err error)
	GetFavorites(ctx context.Context, req *GetFavoritesRequest, opts ...http.CallOption) (rsp *Favorites, err error)
	GetUserProfile(ctx context.Context, req *GetProfileRequest, opts ...http.CallOption) (rsp *GetProfileResponse, err error)
	GetUsers(ctx context.Context, req *GetUsersRequest, opts ...http.CallOption) (rsp *GetUsersResponse, err error)
	ListCreditCards(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *CreditCards, err error)
	SetFavorites(ctx context.Context, req *UpdateFavoritesRequest, opts ...http.CallOption) (rsp *UpdateFavoritesResply, err error)
	UpdateConsumerAddresses(ctx context.Context, req *ConsumerAddress, opts ...http.CallOption) (rsp *ConsumerAddress, err error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...http.CallOption) (rsp *UpdateUserResponse, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) CreateConsumerAddresses(ctx context.Context, in *ConsumerAddress, opts ...http.CallOption) (*ConsumerAddress, error) {
	var out ConsumerAddress
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceCreateConsumerAddresses))
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

func (c *UserServiceHTTPClientImpl) DeleteConsumerAddresses(ctx context.Context, in *DeleteConsumerAddressesRequest, opts ...http.CallOption) (*DeleteConsumerAddressesReply, error) {
	var out DeleteConsumerAddressesReply
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceDeleteConsumerAddresses))
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

func (c *UserServiceHTTPClientImpl) DeleteFavorites(ctx context.Context, in *UpdateFavoritesRequest, opts ...http.CallOption) (*UpdateFavoritesResply, error) {
	var out UpdateFavoritesResply
	pattern := "/v1/users/favorites"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceDeleteFavorites))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...http.CallOption) (*DeleteUserResponse, error) {
	var out DeleteUserResponse
	pattern := "/v1/users"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceDeleteUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetConsumerAddress(ctx context.Context, in *GetConsumerAddressRequest, opts ...http.CallOption) (*ConsumerAddress, error) {
	var out ConsumerAddress
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetConsumerAddress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) GetConsumerAddresses(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*GetConsumerAddressesReply, error) {
	var out GetConsumerAddressesReply
	pattern := "/v1/users/addresses"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetConsumerAddresses))
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

func (c *UserServiceHTTPClientImpl) GetFavorites(ctx context.Context, in *GetFavoritesRequest, opts ...http.CallOption) (*Favorites, error) {
	var out Favorites
	pattern := "/v1/users/favorites"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetFavorites))
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

func (c *UserServiceHTTPClientImpl) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...http.CallOption) (*GetUsersResponse, error) {
	var out GetUsersResponse
	pattern := "/v1/users"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceGetUsers))
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

func (c *UserServiceHTTPClientImpl) SetFavorites(ctx context.Context, in *UpdateFavoritesRequest, opts ...http.CallOption) (*UpdateFavoritesResply, error) {
	var out UpdateFavoritesResply
	pattern := "/v1/users/favorites"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceSetFavorites))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateConsumerAddresses(ctx context.Context, in *ConsumerAddress, opts ...http.CallOption) (*ConsumerAddress, error) {
	var out ConsumerAddress
	pattern := "/v1/users/address"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateConsumerAddresses))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PATCH", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...http.CallOption) (*UpdateUserResponse, error) {
	var out UpdateUserResponse
	pattern := "/v1/users/{user_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
