package service

import (
	v1 "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"backend/pkg/token"
	"context"
	"errors"
)

func (s *UserService) CreateAddresses(ctx context.Context, req *v1.Address) (*v1.Address, error) {
	address, err := s.uc.CreateAddress(ctx, &biz.Address{
		Owner:         req.Owner,
		Name:          req.Name,
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
		Owner:         address.Owner,
		Name:          address.Name,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		StreetAddress: address.StreetAddress,
	}, nil

}
func (s *UserService) UpdateAddresses(ctx context.Context, req *v1.Address) (*v1.Address, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Owner != payload.Owner || req.Name != payload.Name {
		return nil, errors.New("invalid token")
	}

	address, err := s.uc.UpdateAddress(ctx, &biz.Address{
		Id:            req.Id,
		Owner:         payload.Owner,
		Name:          payload.Name,
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
		Owner:         address.Owner,
		Name:          address.Name,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
	}, nil

}
func (s *UserService) DeleteAddresses(ctx context.Context, req *v1.DeleteAddressesRequest) (*v1.DeleteAddressesReply, error) {
	reply, err := s.uc.DeleteAddress(ctx, &biz.DeleteAddressesRequest{
		AddressId: uint32(req.AddressesId),
		Owner:     req.Owner,
		Name:      req.Name,
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
func (s *UserService) GetAddresses(ctx context.Context, req *v1.GetAddressesRequest) (*v1.GetAddressesReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Owner != payload.Owner || req.Name != payload.Name {
		return nil, errors.New("invalid token")
	}

	addresses, err := s.uc.GetAddresses(ctx, &biz.Request{
		Owner: payload.Owner,
		Name:  payload.Name,
	})
	if err != nil {
		return nil, err
	}
	addressList := make([]*v1.Address, len(addresses.Addresses))
	for i, address := range addresses.Addresses {
		addressList[i] = &v1.Address{
			Id:            address.Id,
			Owner:         address.Owner,
			Name:          address.Name,
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
