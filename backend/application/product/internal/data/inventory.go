package data

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/go-kratos/kratos/v2/errors"

	"backend/application/product/internal/data/models"

	"backend/application/product/internal/biz"
)

func (p *productRepo) UpdateInventory(ctx context.Context, req *biz.UpdateInventoryRequest) (*biz.UpdateInventoryReply, error) {
	// 更新库存
	result, err := p.data.db.UpdateInventory(ctx, models.UpdateInventoryParams{
		ProductID:  req.ProductId,
		MerchantID: req.MerchantId,
		Delta:      req.Stock,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		// 库存不足地处理, 即违反了SQL的约束
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23514" {
				return nil, errors.New(400, "INSUFFICIENT_STOCK", "insufficient stock")
			}
			p.log.Debugf("pgErr: %+v, pgErr%+v pgErr%+v", pgErr.Code, pgErr.Message, pgErr.Detail)
			p.log.Debugf("err2: %v", err)
		}
		// if strings.Contains(err.Error(), "violates check constraint") {
		// 	return nil, errors.New(400, "INSUFFICIENT_STOCK", "insufficient stock")
		// }
		return nil, errors.New(500, "INTERNAL_ERROR", "failed to update inventory")
	}

	p.log.Debugf("result: %v", result)
	return &biz.UpdateInventoryReply{
		ProductId:  req.ProductId,
		MerchantId: req.MerchantId,
		Stock:      uint32(req.Stock),
	}, nil
}
