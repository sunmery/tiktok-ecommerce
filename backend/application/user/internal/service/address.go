package service

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/metadata"

	"github.com/google/uuid"

	v1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
)

// CreateAddresses 创建地址
func (s *UserService) CreateAddresses(ctx context.Context, req *v1.Address) (*v1.Address, error) {
	var userStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		md.Get("x-global-user-id")
	}
	userId := uuid.MustParse(userStr)
	address, err := s.uc.CreateAddress(ctx, &biz.Address{
		UserId:        userId,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
	})
	if err != nil {
		return nil, err
	}
	return &v1.Address{
		Id:            address.Id,
		UserId:        userId.String(),
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		StreetAddress: address.StreetAddress,
	}, nil
}

// UpdateAddresses 更新地址
func (s *UserService) UpdateAddresses(ctx context.Context, req *v1.Address) (*v1.Address, error) {
	var userStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		md.Get("x-global-user-id")
	}
	userId := uuid.MustParse(userStr)
	address, err := s.uc.UpdateAddress(ctx, &biz.Address{
		Id:            req.Id,
		UserId:        userId,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
	})
	if err != nil {
		return nil, err
	}
	return &v1.Address{
		Id:            address.Id,
		UserId:        userId.String(),
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
	}, nil
}

// DeleteAddresses 删除地址
func (s *UserService) DeleteAddresses(ctx context.Context, req *v1.DeleteAddressesRequest) (*v1.DeleteAddressesReply, error) {
	var userStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		md.Get("x-global-user-id")
	}
	userId := uuid.MustParse(userStr)

	reply, err := s.uc.DeleteAddress(ctx, &biz.DeleteAddressesRequest{
		AddressId: uint32(req.AddressesId),
		UserId:    userId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DeleteAddressesReply{
		Message: reply.Message,
		Id:      reply.Id,
		Code:    reply.Code,
	}, nil
}

// GetAddresses 获取地址
func (s *UserService) GetAddresses(ctx context.Context, _ *emptypb.Empty) (*v1.GetAddressesReply, error) {
	var userStr string
	if md, ok := metadata.FromServerContext(ctx); ok {
		md.Get("x-global-user-id")
	}
	userId,err := uuid.Parse(userStr)
	if err != nil {
		return nil, errors.New("错误的UserID")
	}
	addresses, err := s.uc.GetAddresses(ctx, &biz.Request{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	addressList := make([]*v1.Address, len(addresses.Addresses))
	for i, address := range addresses.Addresses {
		addressList[i] = &v1.Address{
			Id:            address.Id,
			UserId:        userId.String(),
			City:          address.City,
			State:         address.State,
			Country:       address.Country,
			ZipCode:       address.ZipCode,
			StreetAddress: address.StreetAddress,
		}
	}

	return &v1.GetAddressesReply{
		Addresses: addressList,
	}, nil
}
