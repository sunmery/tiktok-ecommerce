package biz

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	WAREHOUSE    AddressType = "WAREHOUSE"    // 仓库地址
	RETURN       AddressType = "RETURN"       // 退货地址
	STORE        AddressType = "STORE"        // 门店地址
	BILLING      AddressType = "BILLING"      // 财务地址
	HEADQUARTERS AddressType = "HEADQUARTERS" // 总部地址
)

type (
	AddressType string

	MerchantAddressn struct {
		Id            int64
		MerchantId    uuid.UUID
		AddressType   AddressType
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
		AddressType   AddressType
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

	GetMerchantAddressRequestn struct {
		Id         int64
		MerchantId uuid.UUID
	}

	ListMerchantAddressesRequestn struct {
		MerchantId  uuid.UUID
		AddressType AddressType
		OnlyDefault bool
		Page        uint32
		PageSize    uint32
	}

	ListMerchantAddressesResponse struct {
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
	// GetMerchantAddress 获取单个地址详情
	GetMerchantAddress(ctx context.Context, req *GetMerchantAddressRequestn) (*MerchantAddress, error)
	// ListMerchantAddresses 列出商家所有地址（支持按类型过滤）
	ListMerchantAddresses(ctx context.Context, req *ListMerchantAddressesRequestn) (*ListMerchantAddressesResponse, error)
	// SetDefaultAddress 设置默认地址（按地址类型）
	SetDefaultAddress(ctx context.Context, req *SetDefaultAddressRequestn) (*MerchantAddress, error)
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

func (uc *AddressUsecase) GetMerchantAddress(ctx context.Context, req *GetMerchantAddressRequestn) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("GetMerchantAddress: %+v", req)
	return uc.repo.GetMerchantAddress(ctx, req)
}

func (uc *AddressUsecase) ListMerchantAddresses(ctx context.Context, req *ListMerchantAddressesRequestn) (*ListMerchantAddressesResponse, error) {
	uc.log.WithContext(ctx).Debugf("ListMerchantAddresses: %+v", req)
	return uc.repo.ListMerchantAddresses(ctx, req)
}

func (uc *AddressUsecase) SetDefaultAddress(ctx context.Context, req *SetDefaultAddressRequestn) (*MerchantAddress, error) {
	uc.log.WithContext(ctx).Debugf("SetDefaultAddress: %+v", req)
	return uc.repo.SetDefaultAddress(ctx, req)
}
