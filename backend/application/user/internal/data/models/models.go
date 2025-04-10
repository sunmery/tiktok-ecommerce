// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UsersAddresses struct {
	ID            int32     `json:"id"`
	UserID        uuid.UUID `json:"userID"`
	StreetAddress string    `json:"streetAddress"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zipCode"`
}

type UsersCreditCards struct {
	ID          int32     `json:"id"`
	UserID      uuid.UUID `json:"userID"`
	Number      string    `json:"number"`
	Currency    string    `json:"currency"`
	Cvv         string    `json:"cvv"`
	ExpYear     string    `json:"expYear"`
	ExpMonth    string    `json:"expMonth"`
	Owner       string    `json:"owner"`
	Name        *string   `json:"name"`
	Type        string    `json:"type"`
	Brand       string    `json:"brand"`
	Country     string    `json:"country"`
	CreatedTime time.Time `json:"createdTime"`
}

type UsersFavorites struct {
	UserID    uuid.UUID          `json:"userID"`
	ProductID uuid.UUID          `json:"productID"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}
