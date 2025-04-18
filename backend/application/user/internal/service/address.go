package service

import (
	"context"
	"fmt"

	"backend/pkg"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
)

// CreateAddresses 创建地址
func (s *UserService) CreateAddresses(ctx context.Context, req *v1.Address) (*v1.Address, error) {
	userId, uErr := pkg.GetMetadataUesrID(ctx)
	if uErr != nil {
		return nil, fmt.Errorf("错误的UserID")
	}
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
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("错误的UserID")
	}
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
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("错误的UserID")
	}

	reply, err := s.uc.DeleteAddress(ctx, &biz.AddressRequest{
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

// GetAddress 根据ID获取地址
func (s *UserService) GetAddress(ctx context.Context, req *v1.GetAddressRequest) (*v1.Address, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("错误的UserID")
	}
	addresses, err := s.uc.GetAddress(ctx, &biz.AddressRequest{
		AddressId: req.AddressId,
		UserId:    userId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Address{
		Id:            addresses.Id,
		UserId:        addresses.UserId.String(),
		City:          addresses.City,
		State:         addresses.State,
		Country:       addresses.Country,
		ZipCode:       addresses.ZipCode,
		StreetAddress: addresses.StreetAddress,
	}, nil
}

// GetAddresses 获取地址列表
func (s *UserService) GetAddresses(ctx context.Context, _ *emptypb.Empty) (*v1.GetAddressesReply, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, fmt.Errorf("错误的UserID")
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
