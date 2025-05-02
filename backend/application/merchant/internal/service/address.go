package service

import (
	"context"

	"github.com/google/uuid"

	"backend/constants"

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
	addressv1.UnimplementedMerchantAddressesServer
	ac *biz.AddressUsecase
}

func (s *AddressService) CreateMerchantAddress(ctx context.Context, req *pb.MerchantAddress) (*pb.MerchantAddress, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}
	snowflakeID := id.SnowflakeID()

	addressType := convertAddressType(req.AddressType)
	address, err := s.ac.CreateMerchantAddress(ctx, &biz.MerchantAddressn{
		Id:            snowflakeID,
		MerchantId:    merchantId,
		AddressType:   addressType,
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

	return &pb.MerchantAddress{
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
		Remarks:       address.Remarks,
		// Latitude:      address.Latitude,
		// Longitude:     address.Longitude,
	}, nil
}

func (s *AddressService) BatchCreateMerchantAddresses(ctx context.Context, req *pb.BatchCreateMerchantAddressesRequest) (*pb.BatchCreateMerchantAddressesReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	addresses := make([]*biz.MerchantAddressn, len(req.Addresses))
	for i, addr := range req.Addresses {
		addresses[i] = &biz.MerchantAddressn{
			Id:            id.SnowflakeID(),
			MerchantId:    merchantId,
			AddressType:   constants.AddressType(addr.AddressType),
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

	pbAddresses := make([]*pb.MerchantAddress, len(result.Addresses))
	for i, addr := range result.Addresses {
		pbAddresses[i] = &pb.MerchantAddress{
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

	return &pb.BatchCreateMerchantAddressesReply{
		SuccessCount: int32(len(result.Addresses)),
		FailedItems:  nil,
	}, nil
}

func (s *AddressService) UpdateMerchantAddress(ctx context.Context, req *pb.MerchantAddress) (*pb.MerchantAddress, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	address, err := s.ac.UpdateMerchantAddress(ctx, &biz.MerchantAddressn{
		Id:            req.Id,
		MerchantId:    merchantId,
		AddressType:   convertAddressType(req.AddressType),
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

	return &pb.MerchantAddress{
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

func (s *AddressService) DeletMerchanteAddress(ctx context.Context, req *pb.DeletMerchanteAddressRequest) (*empty.Empty, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	return s.ac.DeleteMerchantAddress(ctx, &biz.DeleteMerchantAddressRequestn{
		Id:         req.Id,
		MerchantId: merchantId,
	})
}

func (s *AddressService) ListFilterAddresses(ctx context.Context, req *pb.ListFilterAddressesRequest) (*pb.ListAddressesReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	result, err := s.ac.ListFilterAddresses(ctx, &biz.ListFilterAddressesRequestn{
		MerchantId:  merchantId,
		AddressType: convertAddressType(*req.AddressType),
		Page:        req.Page,
		PageSize:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	addresses := make([]*pb.MerchantAddress, len(result.Addresses))
	for i, addr := range result.Addresses {
		addresses[i] = &pb.MerchantAddress{
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

	return &pb.ListAddressesReply{
		Addresses:  addresses,
		TotalCount: uint32(len(addresses)),
	}, nil
}

func (s *AddressService) GetDefaultAddress(ctx context.Context, req *pb.GetDefaultAddressRequest) (*pb.MerchantAddress, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	address, err := s.ac.GetDefaultAddress(ctx, &biz.GetDefaultAddressRequest{
		MerchantId:  merchantId,
		AddressType: convertAddressType(req.AddressType),
	})
	if err != nil {
		return nil, err
	}

	return convertToPBMerchantAddress(address), nil
}

func (s *AddressService) GetMerchantAddress(ctx context.Context, req *pb.GetMerchantAddressRequest) (*pb.MerchantAddress, error) {
	var merchantId uuid.UUID
	var err error
	if req.MerchantId == nil {
		merchantId, err = pkg.GetMetadataUesrID(ctx)
		if err != nil {
			return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
		}
	} else {
		merchantId, err = uuid.Parse(*req.MerchantId)
		if err != nil {
			return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
		}
	}

	address, err := s.ac.GetMerchantAddress(ctx, &biz.GetMerchantAddressRequest{
		Id:         req.Id,
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	return convertToPBMerchantAddress(address), nil
}

func (s *AddressService) GetDefaultAddresses(ctx context.Context, req *pb.GetDefaultAddressesRequest) (*pb.ListAddressesReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	result, err := s.ac.GetDefaultAddresses(ctx, &biz.GetDefaultAddressesRequest{
		MerchantId: merchantId,
	})
	if err != nil {
		return nil, err
	}

	addresses := make([]*pb.MerchantAddress, len(result.Addresses))
	for i, addr := range result.Addresses {
		addresses[i] = convertToPBMerchantAddress(addr)
	}

	return &pb.ListAddressesReply{
		Addresses:  addresses,
		TotalCount: uint32(len(addresses)),
	}, nil
}

// 辅助函数用于转换为protobuf消息
func convertToPBMerchantAddress(address *biz.MerchantAddress) *pb.MerchantAddress {
	return &pb.MerchantAddress{
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
		Remarks:       address.Remarks,
	}
}

func (s *AddressService) ListAddresses(ctx context.Context, req *pb.ListAddressesRequest) (*pb.ListAddressesReply, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, errors.New(400, "INVALID_MERCHANT_ID", "invalid merchant id")
	}

	result, err := s.ac.ListAddresses(ctx, &biz.ListAddressesRequest{
		MerchantId: merchantId,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	addresses := make([]*pb.MerchantAddress, len(result.Addresses))
	for i, addr := range result.Addresses {
		addresses[i] = &pb.MerchantAddress{
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

	return &pb.ListAddressesReply{
		Addresses:  addresses,
		TotalCount: uint32(len(addresses)),
	}, nil
}

func (s *AddressService) SetDefaultMerchantAddress(ctx context.Context, req *pb.SetDefaultMerchantAddressRequest) (*pb.MerchantAddress, error) {
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

	return &pb.MerchantAddress{
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

func convertAddressType(AddressType pb.MerchantAddress_MerchantAddressType) constants.AddressType {
	switch AddressType {
	case 0:
		return constants.WAREHOUSE
	case 1:
		return constants.RETURN
	case 2:
		return constants.STORE
	case 3:
		return constants.BILLING
	case 4:
		return constants.HEADQUARTERS
	default:
		return constants.WAREHOUSE
	}
}

func convertToPBAddressType(AddressType constants.AddressType) pb.MerchantAddress_MerchantAddressType {
	switch AddressType {
	case constants.WAREHOUSE:
		return pb.MerchantAddress_WAREHOUSE
	case constants.RETURN:
		return pb.MerchantAddress_RETURN
	case constants.STORE:
		return pb.MerchantAddress_STORE
	case constants.BILLING:
		return pb.MerchantAddress_BILLING
	case constants.HEADQUARTERS:
		return pb.MerchantAddress_HEADQUARTERS
	default:
		return pb.MerchantAddress_WAREHOUSE
	}
}
