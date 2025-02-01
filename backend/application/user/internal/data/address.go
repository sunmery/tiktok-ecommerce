package data

import (
	"backend/application/user/internal/biz"
	"backend/application/user/internal/data/models"
	"backend/pkg/token"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func (u *userRepo) CreateAddress(ctx context.Context, req *biz.Address) (*biz.Address, error) {
	address, err := u.data.db.CreatAddress(ctx, models.CreatAddressParams{
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

	return &biz.Address{
		Id:            uint32(address.ID),
		Owner:         address.Owner,
		Name:          address.Name,
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
		Owner:         req.Owner,
		Name:          req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &biz.Address{
		Id:            uint32(address.ID),
		Owner:         address.Owner,
		Name:          address.Name,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		StreetAddress: address.StreetAddress,
	}, err
}

func (u *userRepo) DeleteAddress(ctx context.Context, req *biz.DeleteAddressesRequest) (*biz.DeleteAddressesReply, error) {
	payload, err := token.ExtractPayload(ctx)
	if err != nil {
		return nil, err
	}

	if req.Owner != payload.Owner || req.Name != payload.Name {
		return nil, errors.New("invalid token")
	}
	reply, err := u.data.db.DeleteAddress(ctx, models.DeleteAddressParams{
		Owner: payload.Owner,
		Name:  payload.Name,
		ID:    int32(req.AddressId),
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

func (u *userRepo) GetAddresses(ctx context.Context, req *biz.Request) (*biz.Addresses, error) {

	addresses, aErr := u.data.db.GetAddresses(ctx, models.GetAddressesParams{
		Owner: req.Owner,
		Name:  req.Name,
	})
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
			Owner:         address.Owner,
			Name:          address.Name,
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
