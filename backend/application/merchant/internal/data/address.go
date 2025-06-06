package data

import (
	"context"
	"errors"

	"backend/constants"

	"backend/application/merchant/internal/pkg/id"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/jackc/pgx/v5"

	"backend/pkg/types"

	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/emptypb"

	"backend/application/merchant/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"

	"backend/application/merchant/internal/biz"
)

type addressRepo struct {
	data *Data
	log  *log.Helper
}

func (a addressRepo) ListAddresses(ctx context.Context, req *biz.ListAddressesRequest) (*biz.ListAddressesResponse, error) {
	addresses, err := a.data.db.ListAddresses(ctx, models.ListAddressesParams{
		MerchantID: req.MerchantId,
		Limit:      int64(req.PageSize),
		Offset:     int64((req.Page - 1) * req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*biz.MerchantAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = &biz.MerchantAddress{
			Id:            addr.ID,
			MerchantId:    addr.MerchantID.String(),
			AddressType:   constants.AddressType(addr.AddressType),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			CreatedAt:     addr.CreatedAt.Time,
			UpdatedAt:     addr.UpdatedAt.Time,
		}
	}

	return &biz.ListAddressesResponse{Addresses: result}, nil
}

func (a addressRepo) BatchCreateAddresses(ctx context.Context, req *biz.BatchCreateAddressesRequestn) (*biz.BatchCreateAddressesResponse, error) {
	ids := make([]int64, len(req.Addresses))
	merchantIds := make([]uuid.UUID, len(req.Addresses))
	addressTypes := make([]string, len(req.Addresses))
	contactPersons := make([]string, len(req.Addresses))
	contactPhones := make([]string, len(req.Addresses))
	streetAddresses := make([]string, len(req.Addresses))
	cities := make([]string, len(req.Addresses))
	states := make([]string, len(req.Addresses))
	countries := make([]string, len(req.Addresses))
	zipCodes := make([]string, len(req.Addresses))
	isDefaults := make([]bool, len(req.Addresses))
	remarks := make([]string, len(req.Addresses))

	// TODO
	// err := (*models.MerchantAddressType).Scan(string(addr.AddressType))
	// if err != nil {
	// 	return nil, kerrors.New(400, "address_type", "invalid address type")
	// }

	for i, addr := range req.Addresses {
		ids[i] = id.SnowflakeID()
		merchantIds[i] = addr.MerchantId
		addressTypes[i] = string(addr.AddressType)
		contactPersons[i] = addr.ContactPerson
		contactPhones[i] = addr.ContactPhone
		streetAddresses[i] = addr.StreetAddress
		cities[i] = addr.City
		states[i] = addr.State
		countries[i] = addr.Country
		zipCodes[i] = addr.ZipCode
		isDefaults[i] = addr.IsDefault
		remarks[i] = addr.Remarks
	}

	params := models.BatchCreateAddressesParams{
		ID:         ids,
		MerchantID: merchantIds,
		// AddressType:   addressTypes, // TODO
		ContactPerson: contactPersons,
		ContactPhone:  contactPhones,
		StreetAddress: streetAddresses,
		City:          cities,
		State:         states,
		Country:       countries,
		ZipCode:       zipCodes,
		IsDefault:     isDefaults,
		Remarks:       remarks,
	}
	a.log.Debugf("params:%+v", params)
	addresses, err := a.data.db.BatchCreateAddresses(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*biz.MerchantAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = &biz.MerchantAddress{
			Id:            addr.ID,
			MerchantId:    addr.MerchantID.String(),
			AddressType:   constants.AddressType(addr.AddressType),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			CreatedAt:     addr.CreatedAt.Time,
			UpdatedAt:     addr.UpdatedAt.Time,
		}
	}

	return &biz.BatchCreateAddressesResponse{Addresses: result}, nil
}

func (a addressRepo) UpdateMerchantAddress(ctx context.Context, req *biz.MerchantAddressn) (*biz.MerchantAddress, error) {
	address, err := a.data.db.UpdateAddress(ctx, models.UpdateAddressParams{
		AddressType:   string(req.AddressType),
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		IsDefault:     req.IsDefault,
		ID:            req.Id,
		MerchantID:    req.MerchantId,
	})
	if err != nil {
		return nil, err
	}

	return &biz.MerchantAddress{
		Id:            address.ID,
		MerchantId:    address.MerchantID.String(),
		AddressType:   constants.AddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt.Time,
		UpdatedAt:     address.UpdatedAt.Time,
	}, nil
}

func (a addressRepo) DeleteMerchantAddress(ctx context.Context, req *biz.DeleteMerchantAddressRequestn) (*emptypb.Empty, error) {
	err := a.data.db.DeleteAddress(ctx, models.DeleteAddressParams{
		ID:         req.Id,
		MerchantID: req.MerchantId,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (a addressRepo) GetDefaultAddress(ctx context.Context, req *biz.GetDefaultAddressRequest) (*biz.MerchantAddress, error) {
	address, err := a.data.db.GetDefaultAddress(ctx, models.GetDefaultAddressParams{
		MerchantID:  req.MerchantId,
		AddressType: string(req.AddressType),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ADDRESS_NOT_FOUND", "default address not found")
		}
		return nil, kerrors.New(500, "INTERNAL_ERROR", "database internal error")
	}

	return convertToMerchantAddress(&address), nil
}

func (a addressRepo) GetMerchantAddress(ctx context.Context, req *biz.GetMerchantAddressRequest) (*biz.MerchantAddress, error) {
	address, err := a.data.db.GetAddress(ctx, models.GetAddressParams{
		MerchantID: req.MerchantId,
		ID:         req.Id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, kerrors.New(404, "ADDRESS_NOT_FOUND", "default address not found")
		}
		return nil, kerrors.New(500, "INTERNAL_ERROR", "database internal error")
	}

	return convertToMerchantAddress(&address), nil
}

func (a addressRepo) GetDefaultAddresses(ctx context.Context, req *biz.GetDefaultAddressesRequest) (*biz.ListAddressesResponse, error) {
	addresses, err := a.data.db.GetDefaultAddresses(ctx, req.MerchantId)
	if err != nil {
		return nil, err
	}

	result := make([]*biz.MerchantAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = convertToMerchantAddress(&addr)
	}

	return &biz.ListAddressesResponse{
		Addresses: result,
		Total:     int64(len(addresses)),
	}, nil
}

// 辅助函数用于转换数据模型
func convertToMerchantAddress(addr *models.MerchantAddresses) *biz.MerchantAddress {
	return &biz.MerchantAddress{
		Id:            addr.ID,
		MerchantId:    addr.MerchantID.String(),
		AddressType:   constants.AddressType(addr.AddressType),
		ContactPerson: addr.ContactPerson,
		ContactPhone:  addr.ContactPhone,
		StreetAddress: addr.StreetAddress,
		City:          addr.City,
		State:         addr.State,
		Country:       addr.Country,
		ZipCode:       addr.ZipCode,
		IsDefault:     addr.IsDefault,
		CreatedAt:     addr.CreatedAt.Time,
		UpdatedAt:     addr.UpdatedAt.Time,
		Remarks:       addr.Remarks,
	}
}

func (a addressRepo) ListFilterAddresses(ctx context.Context, req *biz.ListFilterAddressesRequestn) (*biz.ListAddressesResponse, error) {
	addresses, err := a.data.db.ListFilterAddresses(ctx, models.ListFilterAddressesParams{
		MerchantID:  req.MerchantId,
		AddressType: string(req.AddressType),
		Limit:       int64(req.PageSize),
		Offset:      int64((req.Page - 1) * req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*biz.MerchantAddress, len(addresses))
	for i, addr := range addresses {
		result[i] = &biz.MerchantAddress{
			Id:            addr.ID,
			MerchantId:    addr.MerchantID.String(),
			AddressType:   constants.AddressType(addr.AddressType),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			CreatedAt:     addr.CreatedAt.Time,
			UpdatedAt:     addr.UpdatedAt.Time,
		}
	}

	return &biz.ListAddressesResponse{Addresses: result}, nil
}

func (a addressRepo) SetDefaultAddress(ctx context.Context, req *biz.SetDefaultAddressRequestn) (*biz.MerchantAddress, error) {
	merchantId := types.ToPgUUID(req.MerchantId)
	address, err := a.data.db.SetDefaultAddress(ctx, models.SetDefaultAddressParams{
		MerchantID: merchantId,
		ID:         &req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &biz.MerchantAddress{
		Id:            address.ID,
		MerchantId:    address.MerchantID.String(),
		AddressType:   constants.AddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt.Time,
		UpdatedAt:     address.UpdatedAt.Time,
	}, nil
}

func (a addressRepo) CreateMerchantAddress(ctx context.Context, req *biz.MerchantAddressn) (*biz.MerchantAddress, error) {
	log.Debugf("req: %+v", req)
	address, err := a.data.db.CreateAddress(ctx, models.CreateAddressParams{
		ID:            req.Id,
		MerchantID:    req.MerchantId,
		AddressType:   string(req.AddressType),
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		IsDefault:     req.IsDefault,
		Remarks:       req.Remarks,
		// Latitude:      req.Latitude,
		// Longitude:     req.Longitude,
	})
	if err != nil {
		return nil, err
	}

	return &biz.MerchantAddress{
		Id:            address.ID,
		MerchantId:    address.MerchantID.String(),
		AddressType:   constants.AddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt.Time,
		UpdatedAt:     address.UpdatedAt.Time,
		Remarks:       address.Remarks,
		// Latitude:      address.Latitude,
		// Longitude:     address.Longitude,
	}, nil
}

func NewAddressRepo(data *Data, logger log.Logger) biz.AddressRepo {
	return &addressRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
