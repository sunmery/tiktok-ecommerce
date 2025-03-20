package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"

	productv1 "backend/api/product/v1"

	"backend/application/cart/internal/biz"
	"backend/application/cart/internal/data/models"

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
	_, err := c.data.db.EmptyCart(ctx, models.EmptyCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return &biz.EmptyCartResp{
				Success: true,
			}, nil
		}
		return nil, err
	}
	return &biz.EmptyCartResp{
		Success: true,
	}, nil
}

// GetCart implements biz.CartRepo.
func (c *cartRepo) GetCart(ctx context.Context, req *biz.GetCartReq) (*biz.GetCartRelpy, error) {
	carts, err := c.data.db.GetCart(ctx, models.GetCartParams{
		UserID:   req.UserId,
		CartName: "cart",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	if len(carts) == 0 {
		fmt.Println("No cart items found")
		return &biz.GetCartRelpy{Items: nil}, nil
	}

	var productIds []string
	var merchantIds []string
	for _, c := range carts {
		productIds = append(productIds, c.ProductID.String())
		merchantIds = append(merchantIds, c.MerchantID.String())
	}

	// 从商品微服务获取商品信息, 例如价格
	products, perr := c.data.productv1.GetProductsBatch(ctx, &productv1.GetProductsBatchRequest{
		ProductIds:  productIds,
		MerchantIds: merchantIds,
	})
	if perr != nil {
		return nil, perr
	}

	var cartItems []*biz.CartInfo
	for _, cart := range carts {
		for _, p := range products.Items {
			productId, err := uuid.Parse(p.Id)
			if err != nil {
				return nil, err
			}
			merchantId, err := uuid.Parse(p.MerchantId)
			if err != nil {
				return nil, err
			}

			if merchantId == cart.MerchantID && productId == cart.ProductID {
				var picture string
				for _, image := range p.Images {
					if image.IsPrimary {
						picture = image.Url
					}
				}
				cartItems = append(cartItems, &biz.CartInfo{
					MerchantId: cart.MerchantID,
					ProductId:  cart.ProductID,
					Quantity:   uint32(cart.Quantity),
					Price:      p.Price,
					Name:       p.Name,
					Picture:    picture,
				})
			}
		}
	}

	return &biz.GetCartRelpy{
		Items: cartItems,
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
	resp, err := c.data.db.UpsertItem(ctx, models.UpsertItemParams{
		UserID:     req.UserId,
		MerchantID: req.MerchantId,
		ProductID:  req.ProductId,
		Quantity:   int32(req.Quantity),
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
