package biz

import (
	"context"

	"github.com/google/uuid"
)

type Request struct {
	UserId uuid.UUID
}

type AddressRequest struct {
	AddressId uint32
	UserId    uuid.UUID
}

type Address struct {
	Id            uint32
	UserId        uuid.UUID
	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       string
}

type Addresses struct {
	Addresses []*Address `json:"addresses"`
}

type DeleteAddressesReply struct {
	Message string
	Id      uint32
	Code    uint32
}

func (cc *UserUsecase) CreateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("CreateAddress: %+v", req)
	return cc.repo.CreateAddress(ctx, req)
}

func (cc *UserUsecase) UpdateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("UpdateAddress: %+v", req)
	return cc.repo.UpdateAddress(ctx, req)
}

func (cc *UserUsecase) DeleteAddress(ctx context.Context, req *AddressRequest) (*DeleteAddressesReply, error) {
	cc.log.WithContext(ctx).Infof("DeleteAddress: %+v", req)
	return cc.repo.DeleteAddress(ctx, req)
}

func (cc *UserUsecase) GetAddress(ctx context.Context, req *AddressRequest) (*Address, error) {
	cc.log.WithContext(ctx).Infof("GetAddresses: %+v", req)
	return cc.repo.GetAddress(ctx, req)
}

func (cc *UserUsecase) GetAddresses(ctx context.Context, req *Request) (*Addresses, error) {
	cc.log.WithContext(ctx).Infof("GetAddresses: %+v", req)
	return cc.repo.GetAddresses(ctx, req)
}
