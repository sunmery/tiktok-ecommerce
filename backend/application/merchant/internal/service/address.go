package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/application/merchant/internal/pkg/id"

	"backend/pkg"

	"github.com/go-kratos/kratos/v2/errors"

	addressv1 "backend/api/merchant/address/v1"
	pb "backend/api/merchant/address/v1"
	"backend/application/merchant/internal/biz"
)

func NewAddressService(ac *biz.AddressUsecase) *AddressService {
	return &AddressService{ac: ac}
}

type AddressService struct {
	addressv1.UnimplementedMerchantAddressServer
	ac *biz.AddressUsecase
}

func (s *AddressService) CreateAddress(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}
	snowflakeID := id.SnowflakeID()

	addressType := convertAddressType(int32(req.AddressType))
	address, err := s.ac.CreateMerchantAddress(ctx, &biz.MerchantAddressn{
		Id:            snowflakeID,
		MerchantId:    merchantId,
		AddressType:   biz.AddressType(addressType),
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		IsDefault:     req.IsDefault,
		// Remarks:       req.Remarks,
		// Latitude:      req.Latitude,
		// Longitude:     req.Longitude,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Address{
		Id:            address.Id,
		MerchantId:    address.MerchantId,
		AddressType:   convertToPBAddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     timestamppb.New(address.CreatedAt),
		UpdatedAt:     timestamppb.New(address.UpdatedAt),
		// Remarks:       address.Remarks,
		// Latitude:      address.Latitude,
		// Longitude:     address.Longitude,
	}, nil
}

func (s *AddressService) BatchCreateAddresses(ctx context.Context, req *pb.BatchCreateAddressesRequest) (*pb.BatchCreateAddressesResponse, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	addresses := make([]*biz.MerchantAddressn, len(req.Addresses))
	for i, addr := range req.Addresses {
		addresses[i] = &biz.MerchantAddressn{
			Id:            id.SnowflakeID(),
			MerchantId:    merchantId,
			AddressType:   biz.AddressType(convertAddressType(int32(addr.AddressType))),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			Remarks:       addr.Remarks,
		}
	}

	result, err := s.ac.BatchCreateAddresses(ctx, &biz.BatchCreateAddressesRequestn{
		MerchantId: merchantId,
		Addresses:  addresses,
	})
	if err != nil {
		return nil, err
	}

	pbAddresses := make([]*pb.Address, len(result.Addresses))
	for i, addr := range result.Addresses {
		pbAddresses[i] = &pb.Address{
			Id:            addr.Id,
			MerchantId:    addr.MerchantId,
			AddressType:   convertToPBAddressType(addr.AddressType),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			CreatedAt:     timestamppb.New(addr.CreatedAt),
			UpdatedAt:     timestamppb.New(addr.UpdatedAt),
		}
	}

	return &pb.BatchCreateAddressesResponse{
		SuccessCount: int32(len(result.Addresses)),
		FailedItems:  nil,
	}, nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, req *pb.Address) (*pb.Address, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	address, err := s.ac.UpdateMerchantAddress(ctx, &biz.MerchantAddressn{
		Id:            req.Id,
		MerchantId:    merchantId,
		AddressType:   biz.AddressType(convertAddressType(int32(req.AddressType))),
		ContactPerson: req.ContactPerson,
		ContactPhone:  req.ContactPhone,
		StreetAddress: req.StreetAddress,
		City:          req.City,
		State:         req.State,
		Country:       req.Country,
		ZipCode:       req.ZipCode,
		IsDefault:     req.IsDefault,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Address{
		Id:            address.Id,
		MerchantId:    address.MerchantId,
		AddressType:   convertToPBAddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     timestamppb.New(address.CreatedAt),
		UpdatedAt:     timestamppb.New(address.UpdatedAt),
	}, nil
}

func (s *AddressService) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*empty.Empty, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	return s.ac.DeleteMerchantAddress(ctx, &biz.DeleteMerchantAddressRequestn{
		Id:         req.Id,
		MerchantId: merchantId,
	})
}

func (s *AddressService) GetAddress(ctx context.Context, req *pb.GetAddressRequest) (*pb.Address, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	address, err := s.ac.GetMerchantAddress(ctx, &biz.GetMerchantAddressRequestn{
		Id:         req.Id,
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Address{
		Id:            address.Id,
		MerchantId:    address.MerchantId,
		AddressType:   convertToPBAddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     timestamppb.New(address.CreatedAt),
		UpdatedAt:     timestamppb.New(address.UpdatedAt),
	}, nil
}

func (s *AddressService) ListAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesResponse, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	result, err := s.ac.ListMerchantAddresses(ctx, &biz.ListMerchantAddressesRequestn{
		MerchantId:  merchantId,
		AddressType: biz.AddressType(convertAddressType(int32(req.AddressType))),
		OnlyDefault: req.OnlyDefault,
		Page:        req.Page,
		PageSize:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	addresses := make([]*pb.Address, len(result.Addresses))
	for i, addr := range result.Addresses {
		addresses[i] = &pb.Address{
			Id:            addr.Id,
			MerchantId:    addr.MerchantId,
			AddressType:   convertToPBAddressType(addr.AddressType),
			ContactPerson: addr.ContactPerson,
			ContactPhone:  addr.ContactPhone,
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			State:         addr.State,
			Country:       addr.Country,
			ZipCode:       addr.ZipCode,
			IsDefault:     addr.IsDefault,
			CreatedAt:     timestamppb.New(addr.CreatedAt),
			UpdatedAt:     timestamppb.New(addr.UpdatedAt),
		}
	}

	return &pb.ListAddressesResponse{
		Addresses:  addresses,
		TotalCount: uint32(result.Total),
	}, nil
}

func (s *AddressService) SetDefaultAddress(ctx context.Context, req *pb.SetDefaultAddressRequest) (*pb.Address, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	address, err := s.ac.SetDefaultAddress(ctx, &biz.SetDefaultAddressRequestn{
		Id:         req.Id,
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Address{
		Id:            address.Id,
		MerchantId:    address.MerchantId,
		AddressType:   convertToPBAddressType(address.AddressType),
		ContactPerson: address.ContactPerson,
		ContactPhone:  address.ContactPhone,
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     timestamppb.New(address.CreatedAt),
		UpdatedAt:     timestamppb.New(address.UpdatedAt),
	}, nil
}

func convertAddressType(AddressType int32) string {
	switch AddressType {
	case 0:
		return string(biz.WAREHOUSE)
	case 1:
		return string(biz.RETURN)
	case 2:
		return string(biz.STORE)
	case 3:
		return string(biz.BILLING)
	case 4:
		return string(biz.HEADQUARTERS)
	default:
		return string(biz.WAREHOUSE)
	}
}

func convertToPBAddressType(AddressType biz.AddressType) pb.Address_AddressType {
	switch AddressType {
	case biz.WAREHOUSE:
		return pb.Address_WAREHOUSE
	case biz.RETURN:
		return pb.Address_RETURN
	case biz.STORE:
		return pb.Address_STORE
	case biz.BILLING:
		return pb.Address_BILLING
	case biz.HEADQUARTERS:
		return pb.Address_HEADQUARTERS
	default:
		return pb.Address_WAREHOUSE
	}
}
