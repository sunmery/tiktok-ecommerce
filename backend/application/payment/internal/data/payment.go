package data

import (
	"backend/application/payment/internal/biz"
	"backend/application/payment/internal/data/models"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentRepo struct {
	data *Data
	log  *log.Helper
}

func NewPaymentRepo(data *Data, logger log.Logger) biz.PaymentRepo {
	return &PaymentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
func (r *PaymentRepo) Create(ctx context.Context, req *biz.CreateRequest) (*biz.CreateReply, error) {
	r.log.Infof("Create request: %+v", req)
	transactionID, _ := uuid.GenerateUUID()

	// 确保 PayAt 字段正确设置
	now := time.Now()
	payAt := pgtype.Timestamptz{
		Time:  now,
		Valid: true,
	}

	record, err := r.data.db.CreatePayRecord(ctx, models.CreatePayRecordParams{
		UserID:        req.UserID,
		OrderID:       req.OrderID,
		TranscationID: transactionID,
		Amount:        req.Amount,
		PayAt:         payAt,
		Status:        "paid",
	})
	if err != nil {
		return nil, err
	}
	return &biz.CreateReply{TransactionID: record.TranscationID}, nil
}
