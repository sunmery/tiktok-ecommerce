package data

import (
	"backend/application/cart/internal/biz"
	"backend/application/cart/internal/data/models"
	"backend/pkg/types"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type cartRepo struct {
	data *Data
	log  *log.Helper
}

func NewCartRepo(data *Data, logger log.Logger) biz.CartRepo {
	return &cartRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// EmptyCart implements biz.CartRepo.
func (c *cartRepo) EmptyCart(ctx context.Context, req *biz.EmptyCartReq) (*biz.EmptyCartResp, error) {
	effected, err := c.data.db.EmptyCart(ctx, models.EmptyCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &biz.EmptyCartResp{
				Success: false,
			}, nil
		}
		return nil, err
	}
	if effected == 0 {
		return &biz.EmptyCartResp{
			Success: false,
		}, nil
	}
	return &biz.EmptyCartResp{
		Success: true,
	}, nil
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartResp, error) {
	cart, err := c.data.db.GetCart(ctx, models.GetCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	c.log.WithContext(ctx).Infof("GetCarterr________________ : %+v", err)
	if err != nil {
		return nil, err
	}
	c.log.WithContext(ctx).Infof("GetCart request : %+v", cart)
	var cartItems []biz.CartItem
	for _, item := range cart {
		price, err := types.NumericToFloat(item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price to float: %v", err)
		}
		var cartitem biz.CartItem
		cartitem.MerchantId = item.MerchantID
		cartitem.ProductId = item.ProductID
		cartitem.Quantity = item.Quantity
		cartitem.Price = price
		cartItems = append(cartItems, cartitem)
	}

	return &biz.GetCartResp{
		Cart: biz.Cart{
			UserId: req.UserId.String(),
			Items:  cartItems,
		},
	}, nil
}

// RemoveCartItem implements biz.CartRepo.
func (c *cartRepo) RemoveCartItem(ctx context.Context, req *biz.RemoveCartItemReq) (*biz.RemoveCartItemResp, error) {
	_, err := c.data.db.RemoveCartItem(ctx, models.RemoveCartItemParams{
		UserID:     req.UserId,
		MerchantID: req.MerchantId,
		ProductID:  req.ProductId,
		CartName:   "cart",
	})
	if err != nil {
		return nil, err
	}
	return &biz.RemoveCartItemResp{
		Success: true,
	}, nil
}

// UpsertItem implements biz.CartRepo.
func (c *cartRepo) UpsertItem(ctx context.Context, req *biz.UpsertItemReq) (*biz.UpsertItemResp, error) {
	price, err := types.Float64ToNumeric(req.Item.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to convert price to numeric: %v", err)
	}

	resp, err := c.data.db.UpsertItem(ctx, models.UpsertItemParams{
		UserID:     req.UserId,
		MerchantID: req.Item.MerchantId,
		ProductID:  req.Item.ProductId,
		Quantity:   req.Item.Quantity,
		Price:      price,
		CartName:   "cart",
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	return &biz.UpsertItemResp{
		Success: true,
	}, nil
}
