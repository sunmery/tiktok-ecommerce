// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: card.sql

package models

import (
	"context"

	"github.com/google/uuid"
)

const DeleteCreditCard = `-- name: DeleteCreditCard :exec
DELETE
FROM users.credit_cards
WHERE id = $1
`

// DeleteCreditCard
//
//	DELETE
//	FROM users.credit_cards
//	WHERE id = $1
func (q *Queries) DeleteCreditCard(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, DeleteCreditCard, id)
	return err
}

const GetCreditCard = `-- name: GetCreditCard :one
SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
FROM users.credit_cards
WHERE user_id = $1
  AND id = $2
`

type GetCreditCardParams struct {
	UserID uuid.UUID `json:"userID"`
	ID     int32     `json:"id"`
}

// GetCreditCard
//
//	SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
//	FROM users.credit_cards
//	WHERE user_id = $1
//	  AND id = $2
func (q *Queries) GetCreditCard(ctx context.Context, arg GetCreditCardParams) (UsersCreditCards, error) {
	row := q.db.QueryRow(ctx, GetCreditCard, arg.UserID, arg.ID)
	var i UsersCreditCards
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Number,
		&i.Currency,
		&i.Cvv,
		&i.ExpYear,
		&i.ExpMonth,
		&i.Owner,
		&i.Name,
		&i.Type,
		&i.Brand,
		&i.Country,
		&i.CreatedTime,
	)
	return i, err
}

const InsertCreditCard = `-- name: InsertCreditCard :exec
INSERT INTO users.credit_cards (user_id, currency, number, cvv, exp_year, exp_month, owner, name, type, brand, country)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

type InsertCreditCardParams struct {
	UserID   uuid.UUID `json:"userID"`
	Currency string    `json:"currency"`
	Number   string    `json:"number"`
	Cvv      string    `json:"cvv"`
	ExpYear  string    `json:"expYear"`
	ExpMonth string    `json:"expMonth"`
	Owner    string    `json:"owner"`
	Name     *string   `json:"name"`
	Type     string    `json:"type"`
	Brand    string    `json:"brand"`
	Country  string    `json:"country"`
}

// InsertCreditCard
//
//	INSERT INTO users.credit_cards (user_id, currency, number, cvv, exp_year, exp_month, owner, name, type, brand, country)
//	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
func (q *Queries) InsertCreditCard(ctx context.Context, arg InsertCreditCardParams) error {
	_, err := q.db.Exec(ctx, InsertCreditCard,
		arg.UserID,
		arg.Currency,
		arg.Number,
		arg.Cvv,
		arg.ExpYear,
		arg.ExpMonth,
		arg.Owner,
		arg.Name,
		arg.Type,
		arg.Brand,
		arg.Country,
	)
	return err
}

const ListCreditCards = `-- name: ListCreditCards :many
SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
FROM users.credit_cards
WHERE user_id = $1
`

// ListCreditCards
//
//	SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
//	FROM users.credit_cards
//	WHERE user_id = $1
func (q *Queries) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]UsersCreditCards, error) {
	rows, err := q.db.Query(ctx, ListCreditCards, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersCreditCards
	for rows.Next() {
		var i UsersCreditCards
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Number,
			&i.Currency,
			&i.Cvv,
			&i.ExpYear,
			&i.ExpMonth,
			&i.Owner,
			&i.Name,
			&i.Type,
			&i.Brand,
			&i.Country,
			&i.CreatedTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
