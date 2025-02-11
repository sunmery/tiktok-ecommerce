package biz

import "context"

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

func (cc *UserUsecase) CreateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("CreateAddress: %+v", req)
	return cc.repo.CreateAddress(ctx, req)
}

func (cc *UserUsecase) UpdateAddress(ctx context.Context, req *Address) (*Address, error) {
	cc.log.WithContext(ctx).Infof("UpdateAddress: %+v", req)
	return cc.repo.UpdateAddress(ctx, req)
}

func (cc *UserUsecase) DeleteAddress(ctx context.Context, req *DeleteAddressesRequest) (*DeleteAddressesReply, error) {
	cc.log.WithContext(ctx).Infof("DeleteAddress: %+v", req)
	return cc.repo.DeleteAddress(ctx, req)
}

func (cc *UserUsecase) GetAddresses(ctx context.Context, req *Request) (*Addresses, error) {
	cc.log.WithContext(ctx).Infof("GetAddresses: %+v", req)
	return cc.repo.GetAddresses(ctx, req)
}
