package service

import (
	"backend/application/payment/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewPaymentService)

func NewPaymentService(oc *biz.PaymentUsecase) *PaymentService {
	return &PaymentService{
		oc: oc,
	}
}
