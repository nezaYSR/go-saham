package data

import "time"

type InsertUserPayload struct {
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Fullname string `gorm:"size:255;not null" json:"fullname"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	ID           int     `json:"id"`
	Fullname     *string `gorm:"size:255;not null" json:"fullname"`
	FirstOrderId *int    `json:"first_order_id,omitempty"`
}

type InsertOrderItemPayload struct {
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Price     int       `json:"price"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

type UpdateOrderItemPayload struct {
	ID        int       `json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Price     int       `json:"price"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

type SigninPayload struct {
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255" json:"password"`
}

type JWTTokenPayload struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

type InsertOrderHistoryPayload struct {
	UserId       int     `json:"user_id"`
	OrderItemId  int     `json:"order_item_id"`
	Descriptions *string `json:"descriptions,omitempty"`
}
