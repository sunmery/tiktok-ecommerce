package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx 通用的事务方法, 通过外部传递函数作为事务的运行内容
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 开始一个事务, 如sql的begin
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// 事务的运行内容
	query := New(tx)
	err = fn(query)
	// 如果事务发生错误
	if err != nil {
		// 如果回滚发生错误, 那么合并两个错误为一个错误返回回去
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err is: '%s', rellback err is: '%s'", err, rbErr)
		}
		// 返回事务错误
		return err
	}

	// 没有错误则提交该事务
	return tx.Commit(ctx)
}
