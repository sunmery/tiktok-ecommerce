package biz

import (
	"github.com/google/wire"

	"github.com/go-kratos/kratos/v2/log"
)

var ProviderSet = wire.NewSet(NewProductUsecase)

// ProductUsecase is a Product usecase.
type ProductUsecase struct {
	repo ProductRepo
	log  *log.Helper
}

// NewProductUsecase new a Product usecase.
func NewProductUsecase(repo ProductRepo, logger log.Logger) *ProductUsecase {
	return &ProductUsecase{repo: repo, log: log.NewHelper(logger)}
}
