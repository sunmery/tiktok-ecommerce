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
	record, err := r.data.db.CreatePayRecord(ctx, models.CreatePayRecordParams{
		UserID:        req.UserID,
		OrderID:       req.OrderID,
		TranscationID: transactionID,
		Amount:        req.Amount,
		PayAt:         pgtype.Timestamptz{Time: time.Now()},
		Status:        "paid",
	})
	if err != nil {
		return nil, err
	}
	return &biz.CreateReply{TransactionID: record.TranscationID}, nil
}
