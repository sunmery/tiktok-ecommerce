package service

import (
	pb "backend/api/user/v1"
	"backend/application/user/internal/biz"
	"backend/pkg/token"
	"context"
	"errors"
)
// NewAddressService new a User service.
func NewAddressService(ac *biz.AddressesUsecase) *AddressService {
	return &AddressService{ac: ac}
}

func (a *AddressService) CreateAddresses(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	address, err := a.ac.CreateAddress(ctx, &biz.Address{
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
	return &pb.Address{
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
func (a *AddressService) UpdateAddresses(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Owner != payload.Owner || req.Name != payload.Name {
		return nil, errors.New("invalid token")
	}

	address, err := a.ac.UpdateAddress(ctx, &biz.Address{
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
	return &pb.Address{
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
func (a *AddressService) DeleteAddresses(ctx context.Context, req *pb.DeleteAddressesRequest) (*pb.DeleteAddressesReply, error) {
	reply, err := a.ac.DeleteAddress(ctx, &biz.DeleteAddressesRequest{
		AddressId: uint32(req.AddressesId),
		Owner:     req.Owner,
		Name:      req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAddressesReply{
		Message: reply.Message,
		Id:      reply.Id,
		Code:    reply.Code,
	}, nil
}
func (a *AddressService) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Owner != payload.Owner || req.Name != payload.Name {
		return nil, errors.New("invalid token")
	}

	addresses, err := a.ac.GetAddresses(ctx, &biz.Request{
		Owner: payload.Owner,
		Name:  payload.Name,
	})
	if err != nil {
		return nil, err
	}
	addressList := make([]*pb.Address, len(addresses.Addresses))
	for i, address := range addresses.Addresses {
		addressList[i] = &pb.Address{
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

	return &pb.GetAddressesReply{
		Addresses: addressList,
	}, nil
}
