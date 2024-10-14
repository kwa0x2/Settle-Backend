package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/kwa0x2/Settle-Backend/types"
	"time"
)

type User struct {
	ID            string         `bson:"_id" validate:"required"` // Steam ID
	Name          string         `bson:"name"  validate:"required,min=2,max=32"`
	Avatar        string         `bson:"avatar"  validate:"required"`
	ProfileURL    string         `bson:"profile_url"  validate:"required"`
	TotalPlaytime int            `bson:"total_playtime"  validate:"gte=500"`
	Role          types.UserRole `bson:"role" validate:"required"`
	CreatedAt     time.Time      `bson:"created_at"  validate:"required"`
	UpdatedAt     time.Time      `bson:"updated_at"  validate:"required"`
	DeletedAt     *time.Time     `bson:"deleted_at,omitempty" `
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
