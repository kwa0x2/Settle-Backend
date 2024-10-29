package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID            string     `bson:"_id" validate:"required"` // Steam ID
	Name          string     `bson:"name"  validate:"required,min=2,max=32"`
	Avatar        string     `bson:"avatar"  validate:"required"`
	ProfileURL    string     `bson:"profile_url"  validate:"required"`
	TotalPlaytime int        `bson:"total_playtime"  validate:"gte=500"`
	CreatedAt     time.Time  `bson:"created_at"  validate:"required"`
	UpdatedAt     time.Time  `bson:"updated_at"  validate:"required"`
	DeletedAt     *time.Time `bson:"deleted_at,omitempty" `
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindOne(ctx context.Context, filter bson.D) (User, error)
	GetDatabase() *mongo.Database
}

type UserUsecase interface {
	Create(user *User) error
	CreateAndJoinRoom(user *User, userRoom *UserRoom) error
}
