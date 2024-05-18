// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"database/sql"
	"time"
)

type Bookable struct {
	ID                   int64          `json:"id"`
	LocationID           int64          `json:"location_id"`
	Capacity             int64          `json:"capacity"`
	Description          sql.NullString `json:"description"`
	DefaultBookingLength sql.NullInt64  `json:"default_booking_length"`
	CustomBookingLength  bool           `json:"custom_booking_length"`
	ClearTime            sql.NullInt64  `json:"clear_time"`
	CreatedAt            time.Time      `json:"created_at"`
}

type Booking struct {
	ID         int64         `json:"id"`
	BookableID sql.NullInt64 `json:"bookable_id"`
	BookedBy   int64         `json:"booked_by"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Capacity   int64         `json:"capacity"`
	CreatedAt  time.Time     `json:"created_at"`
}

type BookingLength struct {
	ID                   int64         `json:"id"`
	BookableID           sql.NullInt64 `json:"bookable_id"`
	DefaultBookingLength sql.NullInt64 `json:"default_booking_length"`
	CreatedAt            time.Time     `json:"created_at"`
}

type Gisdatum struct {
	ID  int64       `json:"id"`
	Geo interface{} `json:"geo"`
}

type Location struct {
	ID             int64          `json:"id"`
	OrganisationID sql.NullInt64  `json:"organisation_id"`
	UserID         sql.NullInt64  `json:"user_id"`
	FullName       sql.NullString `json:"full_name"`
	Line1          sql.NullString `json:"line1"`
	Line2          sql.NullString `json:"line2"`
	City           sql.NullString `json:"city"`
	County         sql.NullString `json:"county"`
	CountryCode    sql.NullString `json:"country_code"`
	Geo            interface{}    `json:"geo"`
	CreatedAt      time.Time      `json:"created_at"`
}

type Organisation struct {
	ID           int64          `json:"id"`
	CountryCode  sql.NullInt32  `json:"country_code"`
	MerchantName sql.NullString `json:"merchant_name"`
	CreatedAt    time.Time      `json:"created_at"`
}

type User struct {
	ID                int64          `json:"id"`
	FullName          sql.NullString `json:"full_name"`
	Email             sql.NullString `json:"email"`
	EncryptedPassword sql.NullString `json:"encrypted_password"`
	Mobile            sql.NullString `json:"mobile"`
	OrganisationID    sql.NullInt64  `json:"organisation_id"`
	CreatedAt         time.Time      `json:"created_at"`
}
