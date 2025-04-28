package biz

import (
	"context"
	"time"

	"backend/constants"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	MerchantAddressn struct {
		Id            int64
		MerchantId    uuid.UUID
		AddressType   constants.AddressType
		ContactPerson string
		ContactPhone  string

		StreetAddress string
		City          string
		State         string
		Country       string
		ZipCode       string
		IsDefault     bool
		CreatedAt     time.Time
		UpdatedAt     time.Time
		Remarks       string
		// Latitude      float64
		// Longitude     float64
	}
	MerchantAddress struct {
		Id            int64
		MerchantId    string
		AddressType   constants.AddressType
		ContactPerson string
		ContactPhone  string

		StreetAddress string
		City          string
		State         string
		Country       string
		ZipCode       string
		IsDefault     bool
		CreatedAt     time.Time
		UpdatedAt     time.Time
		Remarks       string
		Latitude      float64
		Longitude     float64
	}
)

type (
	BatchCreateAddressesRequestn struct {
		MerchantId uuid.UUID
		Addresses  []*MerchantAddressn
	}

	BatchCreateAddressesResponse struct {
		Addresses []*MerchantAddress
	}

	DeleteMerchantAddressRequestn struct {
		Id         int64
		MerchantId uuid.UUID
	}
	GetMerchantAddressRequest struct {
		Id         int64
		MerchantId uuid.UUID
	}

	ListFilterAddressesRequestn struct {
		MerchantId  uuid.UUID
		AddressType constants.AddressType

		Page     uint32
		PageSize uint32
	}
	GetDefaultAddressesRequest struct {
		MerchantId uuid.UUID
		Page       uint32
		PageSize   uint32
	}
	GetDefaultAddressRequest struct {
		MerchantId  uuid.UUID
		AddressType constants.AddressType
		Page        uint32
		PageSize    uint32
	}
	ListAddressesRequest struct {
		MerchantId uuid.UUID
		Page       uint32
		PageSize   uint32
	}

	ListAddressesResponse struct {
		Addresses []*MerchantAddress
		Total     int64
	}

	SetDefaultAddressRequestn struct {
		Id         int64
		MerchantId uuid.UUID
	}

	GetShippingAddressRequestn struct {
		MerchantId uuid.UUID
	}
)

func NewAddressUsecase(repo AddressRepo, logger log.Logger) *AddressUsecase {
	return &AddressUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

type AddressUsecase struct {
	repo AddressRepo
	log  *log.Helper
}

// AddressRepo 地址域方法
type AddressRepo interface {
	// CreateMerchantAddress 创建商家地址（支持多类型地址）
	CreateMerchantAddress(ctx context.Context, req *MerchantAddressn) (*MerchantAddress, error)
	// BatchCreateAddresses 批量导入商家地址（CSV/JSON格式）
	BatchCreateAddresses(ctx context.Context, req *BatchCreateAddressesRequestn) (*BatchCreateAddressesResponse, error)
	// UpdateMerchantAddress 更新商家地址（支持部分更新）
	UpdateMerchantAddress(ctx context.Context, req *MerchantAddressn) (*MerchantAddress, error)
	// DeleteMerchantAddress 删除商家地址
	DeleteMerchantAddress(ctx context.Context, req *DeleteMerchantAddressRequestn) (*emptypb.Empty, error)
	// ListFilterAddresses 获取单个地址详情
	ListFilterAddresses(ctx context.Context, req *ListFilterAddressesRequestn) (*ListAddressesResponse, error)
	// ListAddresses 列出全部商家地址
	ListAddresses(ctx context.Context, req *ListAddressesRequest) (*ListAddressesResponse, error)
	// SetDefaultAddress 设置默认地址（按地址类型）
	SetDefaultAddress(ctx context.Context, req *SetDefaultAddressRequestn) (*MerchantAddress, error)
	// GetDefaultAddress 获取默认地址（按地址类型）
	GetDefaultAddress(ctx context.Context, req *GetDefaultAddressRequest) (*MerchantAddress, error)
	// GetDefaultAddresses 获取全部地址类型的默认地址
	GetDefaultAddresses(ctx context.Context, req *GetDefaultAddressesRequest) (*ListAddressesResponse, error)
	// GetMerchantAddress 获取单个地址详情
	GetMerchantAddress(ctx context.Context, req *GetMerchantAddressRequest) (*MerchantAddress, error)
}

func (uc *AddressUsecase) CreateMerchantAddress(ctx context.Context, req *MerchantAddressn) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("CreateMerchantAddress: %+v", req)
	return uc.repo.CreateMerchantAddress(ctx, req)
}

func (uc *AddressUsecase) BatchCreateAddresses(ctx context.Context, req *BatchCreateAddressesRequestn) (*BatchCreateAddressesResponse, error) {
	uc.log.WithContext(ctx).Debugf("BatchCreateAddresses: %+v", req)
	return uc.repo.BatchCreateAddresses(ctx, req)
}

func (uc *AddressUsecase) UpdateMerchantAddress(ctx context.Context, req *MerchantAddressn) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("UpdateMerchantAddress: %+v", req)
	return uc.repo.UpdateMerchantAddress(ctx, req)
}

func (uc *AddressUsecase) DeleteMerchantAddress(ctx context.Context, req *DeleteMerchantAddressRequestn) (*emptypb.Empty, error) {
	uc.log.WithContext(ctx).Debugf("DeleteMerchantAddress: %+v", req)
	return uc.repo.DeleteMerchantAddress(ctx, req)
}

func (uc *AddressUsecase) ListFilterAddresses(ctx context.Context, req *ListFilterAddressesRequestn) (*ListAddressesResponse, error) {
	uc.log.WithContext(ctx).Debugf("ListFilterAddresses: %+v", req)
	return uc.repo.ListFilterAddresses(ctx, req)
}

func (uc *AddressUsecase) ListAddresses(ctx context.Context, req *ListAddressesRequest) (*ListAddressesResponse, error) {
	uc.log.WithContext(ctx).Debugf("ListFilterAddresses: %+v", req)
	return uc.repo.ListAddresses(ctx, req)
}

func (uc *AddressUsecase) GetDefaultAddress(ctx context.Context, req *GetDefaultAddressRequest) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("GetDefaultAddress: %+v", req)
	return uc.repo.GetDefaultAddress(ctx, req)
}

func (uc *AddressUsecase) GetDefaultAddresses(ctx context.Context, req *GetDefaultAddressesRequest) (*ListAddressesResponse, error) {
	uc.log.WithContext(ctx).Debugf("GetDefaultAddresses: %+v", req)
	return uc.repo.GetDefaultAddresses(ctx, req)
}

func (uc *AddressUsecase) GetMerchantAddress(ctx context.Context, req *GetMerchantAddressRequest) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("GetMerchantAddress: %+v", req)
	return uc.repo.GetMerchantAddress(ctx, req)
}

func (uc *AddressUsecase) SetDefaultAddress(ctx context.Context, req *SetDefaultAddressRequestn) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("SetDefaultAddress: %+v", req)
	return uc.repo.SetDefaultAddress(ctx, req)
}
