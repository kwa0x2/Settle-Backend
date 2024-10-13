package models

import "time"

type User struct {
	ID            string     `bson:"_id" validate:"required"` // Steam ID
	Name          string     `bson:"name"  validate:"required,min=2,max=32"`
	Avatar        string     `bson:"avatar"  validate:"required"`
	ProfileURL    string     `bson:"profile_url"  validate:"required"`
	TotalPlaytime int        `bson:"total_playtime"  validate:"gte=500"`
	CreatedAt     time.Time  `bson:"created_at" `
	UpdatedAt     time.Time  `bson:"updated_at"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" `
}
