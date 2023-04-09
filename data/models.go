package data

import (
	"database/sql"
	"time"
)

var Db *sql.DB

func NullableInt(value int, condition bool) sql.NullInt64 {
	if condition == false {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(value), Valid: true}
}

type NullableTime struct {
	Time  time.Time
	Valid bool
}

type Users struct {
	ID           int           `gorm:"primary_key;auto_increment" json:"id"`
	Username     string        `gorm:"size:255;not null;unique" json:"username"`
	Fullname     string        `gorm:"size:255;not null" json:"fullname"`
	FirstOrderId *int          `json:"first_order_id,omitempty"`
	Password     string        `gorm:"password" json:"password"`
	Role         UserRole      `gorm:"size:100;not null;" json:"role"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *NullableTime `json:"updated_at,omitempty"`
	DeletedAt    *NullableTime `json:"deleted_at,omitempty"`
}

type OrdersItem struct {
	ID        int           `gorm:"primary_key;auto_increment" json:"id"`
	Name      string        `gorm:"size:255;not null;unique" json:"name"`
	Price     int           `json:"price"`
	ExpiredAt time.Time     `json:"expired_at"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *NullableTime `json:"updated_at,omitempty"`
	DeletedAt *NullableTime `json:"deleted_at,omitempty"`
}

type OrdersHistories struct {
	ID           int           `gorm:"primary_key;auto_increment" json:"id"`
	UserId       int           `json:"user_id"`
	OrderItemId  int           `json:"order_item_id"`
	Descriptions *string       `gorm:"size:255" json:"descriptions"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *NullableTime `json:"updated_at,omitempty"`
	DeletedAt    *NullableTime `json:"deleted_at,omitempty"`
}
