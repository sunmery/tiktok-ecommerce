// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	//CreatAddress
	//
	//  INSERT INTO users.addresses(user_id, street_address, city, state, country, zip_code)
	//  VALUES ($1, $2, $3, $4, $5, $6)
	//  RETURNING id, user_id, street_address, city, state, country, zip_code
	CreatAddress(ctx context.Context, arg CreatAddressParams) (UsersAddresses, error)
	//DeleteAddress
	//
	//  DELETE
	//  FROM users.addresses
	//  WHERE id = $1
	//    AND user_id = $2
	//  RETURNING id, user_id, street_address, city, state, country, zip_code
	DeleteAddress(ctx context.Context, arg DeleteAddressParams) (UsersAddresses, error)
	//DeleteCreditCard
	//
	//  DELETE
	//  FROM users.credit_cards
	//  WHERE id = $1
	DeleteCreditCard(ctx context.Context, id int32) error
	//GetAddress
	//
	//  SELECT id, user_id, street_address, city, state, country, zip_code
	//  FROM users.addresses
	//  WHERE id = $1
	//    AND user_id = $2
	//  LIMIT 1
	GetAddress(ctx context.Context, arg GetAddressParams) (UsersAddresses, error)
	//GetAddresses
	//
	//  SELECT id, user_id, street_address, city, state, country, zip_code
	//  FROM users.addresses
	//  WHERE user_id = $1
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]UsersAddresses, error)
	//GetCreditCard
	//
	//  SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
	//  FROM users.credit_cards
	//  WHERE user_id = $1
	//    AND id = $2
	GetCreditCard(ctx context.Context, arg GetCreditCardParams) (UsersCreditCards, error)
	//InsertCreditCard
	//
	//  INSERT INTO users.credit_cards (user_id, currency, number, cvv, exp_year, exp_month, owner, name, type, brand, country)
	//  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	InsertCreditCard(ctx context.Context, arg InsertCreditCardParams) error
	//ListCreditCards
	//
	//  SELECT id, user_id, number, currency, cvv, exp_year, exp_month, owner, name, type, brand, country, created_time
	//  FROM users.credit_cards
	//  WHERE user_id = $1
	ListCreditCards(ctx context.Context, userID uuid.UUID) ([]UsersCreditCards, error)
	//UpdateAddress
	//
	//  UPDATE users.addresses
	//  SET street_address = coalesce($1, street_address),
	//      city           = coalesce($2, city),
	//      state          = coalesce($3, state),
	//      country        = coalesce($4, country),
	//      zip_code       = coalesce($5, zip_code)
	//  WHERE id = $6
	//    AND user_id = $7
	//  RETURNING id, user_id, street_address, city, state, country, zip_code
	UpdateAddress(ctx context.Context, arg UpdateAddressParams) (UsersAddresses, error)
}

var _ Querier = (*Queries)(nil)
