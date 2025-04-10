// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	//CreateOrder
	//
	//  INSERT INTO orders.orders (id, user_id, currency, street_address,
	//                             city, state, country, zip_code, email)
	//  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	//  RETURNING id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
	CreateOrder(ctx context.Context, arg CreateOrderParams) (OrdersOrders, error)
	//CreateSubOrder
	//
	//  INSERT INTO orders.sub_orders (id, order_id, merchant_id, total_amount,
	//                                 currency, status, items)
	//  VALUES ($1, $2, $3, $4, $5, $6, $7)
	//  RETURNING id, order_id, merchant_id, total_amount, currency, status, items, created_at, updated_at, payment_status
	CreateSubOrder(ctx context.Context, arg CreateSubOrderParams) (OrdersSubOrders, error)
	//GetConsumerOrders
	//
	//  SELECT o.id, o.user_id, o.currency, o.street_address, o.city, o.state, o.country, o.zip_code, o.email, o.created_at, o.updated_at, o.payment_status,
	//         json_agg(
	//                 json_build_object(
	//                         'id', so.id,
	//                         'merchant_id', so.merchant_id,
	//                         'total_amount', so.total_amount,
	//                         'currency', so.currency,
	//                         'status', so.status,
	//                         'items', so.items,
	//                         'created_at', so.created_at,
	//                         'updated_at', so.updated_at
	//                 )
	//         ) AS sub_orders
	//  FROM orders.orders o
	//           LEFT JOIN orders.sub_orders so ON o.id = so.order_id
	//  WHERE o.user_id = $1
	//  GROUP BY o.id
	//  LIMIT $3 OFFSET $2
	GetConsumerOrders(ctx context.Context, arg GetConsumerOrdersParams) ([]GetConsumerOrdersRow, error)
	//GetOrderByID
	//
	//  SELECT o.id, o.user_id, o.currency, o.street_address, o.city, o.state, o.country, o.zip_code, o.email, o.created_at, o.updated_at, o.payment_status,
	//         json_agg(
	//                 json_build_object(
	//                         'id', so.id,
	//                         'merchant_id', so.merchant_id,
	//                         'total_amount', so.total_amount,
	//                         'currency', so.currency,
	//                         'status', so.status,
	//                         'items', so.items,
	//                         'created_at', so.created_at,
	//                         'updated_at', so.updated_at
	//                 )
	//         ) AS sub_orders
	//  FROM orders.orders o
	//           LEFT JOIN orders.sub_orders so ON o.id = so.order_id
	//  WHERE o.user_id = $1
	//    AND o.id = $2
	//  GROUP BY o.id
	GetOrderByID(ctx context.Context, arg GetOrderByIDParams) (GetOrderByIDRow, error)
	//GetOrderByUserID
	//
	//  SELECT id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
	//  FROM orders.orders
	//  WHERE user_id = $1
	GetOrderByUserID(ctx context.Context, userID uuid.UUID) (OrdersOrders, error)
	//GetUserOrdersWithSuborders
	//
	//  SELECT o.id         AS order_id,
	//         o.currency   AS order_currency,
	//         o.street_address,
	//         o.city,
	//         o.state,
	//         o.country,
	//         o.zip_code,
	//         o.email,
	//         o.created_at AS order_created,
	//         jsonb_agg(
	//                 jsonb_build_object(
	//                         'suborder_id', so.id,
	//                         'merchant_id', so.merchant_id,
	//                         'total_amount', so.total_amount,
	//                         'currency', so.currency,
	//                         'status', so.status,
	//                         'items', so.items,
	//                         'created_at', so.created_at,
	//                         'updated_at', so.updated_at
	//                 ) ORDER BY so.created_at
	//         )            AS suborders
	//  FROM orders.orders o
	//           LEFT JOIN orders.sub_orders so ON o.id = so.order_id
	//  WHERE o.user_id = $1::uuid
	//  GROUP BY o.id, o.currency, o.street_address, o.city, o.state, o.country, o.zip_code, o.email, o.created_at
	//  ORDER BY o.created_at DESC
	GetUserOrdersWithSuborders(ctx context.Context, dollar_1 uuid.UUID) ([]GetUserOrdersWithSubordersRow, error)
	//ListOrders
	//
	//  SELECT id, order_id, merchant_id, total_amount, currency, status, items, created_at, updated_at, payment_status
	//  FROM orders.sub_orders
	//  ORDER BY created_at DESC
	//  LIMIT $2 OFFSET $1
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]OrdersSubOrders, error)
	//MarkOrderAsPaid
	//
	//  UPDATE orders.orders
	//  SET payment_status = $1,
	//      updated_at     = now()
	//  WHERE id = $2
	//  RETURNING id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
	MarkOrderAsPaid(ctx context.Context, arg MarkOrderAsPaidParams) (OrdersOrders, error)
	//MarkSubOrderAsPaid
	//
	//  UPDATE orders.sub_orders
	//  SET payment_status = $1,
	//      updated_at     = now()
	//  WHERE order_id = $2
	//  RETURNING id, order_id, merchant_id, total_amount, currency, status, items, created_at, updated_at, payment_status
	MarkSubOrderAsPaid(ctx context.Context, arg MarkSubOrderAsPaidParams) (OrdersSubOrders, error)
	//QuerySubOrders
	//
	//  SELECT id,
	//         merchant_id,
	//         total_amount,
	//         currency,
	//         status,
	//         items,
	//         created_at,
	//         updated_at
	//  FROM orders.sub_orders
	//  WHERE order_id = $1
	//  ORDER BY created_at
	QuerySubOrders(ctx context.Context, orderID int64) ([]QuerySubOrdersRow, error)
	//UpdateOrderPaymentStatus
	//
	//  UPDATE orders.orders
	//  SET payment_status = $2,
	//      updated_at     = now()
	//  WHERE id = $1
	UpdateOrderPaymentStatus(ctx context.Context, arg UpdateOrderPaymentStatusParams) error
	//UpdatePaymentStatus
	//
	//  SELECT id, user_id, payment_status
	//  FROM orders.orders
	//  WHERE id = $1
	//      FOR UPDATE
	UpdatePaymentStatus(ctx context.Context, id int64) (UpdatePaymentStatusRow, error)
	//UpdateSubOrderStatus
	//
	//  UPDATE orders.sub_orders
	//  SET status     = $2,
	//      updated_at = $3
	//  WHERE id = $1
	UpdateSubOrderStatus(ctx context.Context, arg UpdateSubOrderStatusParams) error
}

var _ Querier = (*Queries)(nil)
