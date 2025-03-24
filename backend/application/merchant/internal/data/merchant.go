package data

import (
	"backend/application/merchant/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type inventoryRepo struct {
	data *Data
	log  *log.Helper
}
type productRepo struct {
	data *Data
	log  *log.Helper
}

func NewInventoryRepo(data *Data, logger log.Logger) biz.InventoryRepo {
	return &inventoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
