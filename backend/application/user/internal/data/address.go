package data

import (
	"context"
	"errors"
	"net/http"

	"backend/application/user/internal/biz"
	"backend/application/user/internal/data/models"

	"github.com/jackc/pgx/v5"
)

func (u *userRepo) CreateAddress(ctx context.Context, req *biz.Address) (*biz.Address, error) {
	address, err := u.data.db.CreatAddress(ctx, models.CreatAddressParams{
		UserID:        req.UserId,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
	})
	if err != nil {
		return nil, err
	}

	return &biz.Address{
		Id:            uint32(address.ID),
		UserId:        req.UserId,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		StreetAddress: address.StreetAddress,
	}, nil
}

func (u *userRepo) UpdateAddress(ctx context.Context, req *biz.Address) (*biz.Address, error) {
	address, err := u.data.db.UpdateAddress(ctx, models.UpdateAddressParams{
		StreetAddress: &req.StreetAddress,
		City:          &req.City,
		State:         &req.State,
		Country:       &req.Country,
		ZipCode:       &req.ZipCode,
		ID:            int32(req.Id),
		UserID:        req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &biz.Address{
		Id:            uint32(address.ID),
		UserId:        address.UserID,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		StreetAddress: address.StreetAddress,
	}, err
}

func (u *userRepo) DeleteAddress(ctx context.Context, req *biz.AddressRequest) (*biz.DeleteAddressesReply, error) {
	reply, err := u.data.db.DeleteAddress(ctx, models.DeleteAddressParams{
		UserID: req.UserId,
		ID:     int32(req.AddressId),
	})
	if err != nil {
		return nil, err
	}
	return &biz.DeleteAddressesReply{
		Message: "OK",
		Id:      uint32(reply.ID),
		Code:    http.StatusOK,
	}, nil
}

func (u *userRepo) GetAddress(ctx context.Context, req *biz.AddressRequest) (*biz.Address, error) {
	address, aErr := u.data.db.GetAddress(ctx, models.GetAddressParams{
		ID:     int32(req.AddressId),
		UserID: req.UserId,
	})
	if aErr != nil {
		if errors.Is(aErr, pgx.ErrNoRows) {
			return &biz.Address{}, nil // 返回空列表
		}
		return nil, aErr
	}

	return &biz.Address{
		Id:            uint32(address.ID),
		UserId:        address.UserID,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
	}, nil
}

func (u *userRepo) GetAddresses(ctx context.Context, req *biz.Request) (*biz.Addresses, error) {
	addresses, aErr := u.data.db.GetAddresses(ctx, req.UserId)
	if aErr != nil {
		if errors.Is(aErr, pgx.ErrNoRows) {
			return &biz.Addresses{Addresses: []*biz.Address{}}, nil // 返回空列表
		}
		return nil, aErr
	}

	addressList := make([]*biz.Address, len(addresses))
	for i, address := range addresses {
		addressList[i] = &biz.Address{
			Id:            uint32(address.ID),
			UserId:        address.UserID,
			StreetAddress: address.StreetAddress,
			City:          address.City,
			State:         address.State,
			Country:       address.Country,
			ZipCode:       address.ZipCode,
		}
	}

	return &biz.Addresses{
		Addresses: addressList,
	}, nil
}
