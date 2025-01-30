package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Request struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type DeleteAddressesRequest struct {
	AddressId uint32 `json:"address_id"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
}

type Address struct {
	Id            uint32 `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}

type Addresses struct {
	Addresses []*Address `json:"addresses"`
}

type DeleteAddressesReply struct {
	Message string `json:"message"`
	Id      uint32 `json:"id"`
	Code    uint32 `json:"code"`
}

// AddressesRepo is a Greater repo.
type AddressesRepo interface {
	CreateAddress(ctx context.Context, req *Address) (*Address, error)
	UpdateAddress(ctx context.Context, req *Address) (*Address, error)
	DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	GetAddresses(ctx context.Context, req *Request) (*Addresses, error)
}

// AddressesUsecase is a Addresses usecase.
type AddressesUsecase struct {
	repo AddressesRepo
	log  *log.Helper
}

// NewAddressesUsecase new a Addresses usecase.
func NewAddressesUsecase(repo AddressesRepo, logger log.Logger) *AddressesUsecase {
	return &AddressesUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateAddress creates a Addresses, and returns the new Addresses.
func (uc *AddressesUsecase) CreateAddress(ctx context.Context, req *Address) (*Address, error) {
	uc.log.WithContext(ctx).Infof("CreateAddress: %+v", req)
	return uc.repo.CreateAddress(ctx, req)
}
func (uc *AddressesUsecase) UpdateAddress(ctx context.Context, req *Address) (*Address, error) {
	uc.log.WithContext(ctx).Infof("UpdateAddress: %+v", req)
	return uc.repo.UpdateAddress(ctx, req)
}
func (uc *AddressesUsecase) DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error) {
	uc.log.WithContext(ctx).Infof("DeleteAddress: %+v", req)
	return uc.repo.DeleteAddress(ctx, req)
}
func (uc *AddressesUsecase) GetAddresses(ctx context.Context, req *Request) (*Addresses, error) {
	uc.log.WithContext(ctx).Infof("GetAddresses: %+v", req)
	return uc.repo.GetAddresses(ctx, req)
}
