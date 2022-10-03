package models

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	Bio       string             `bson:"bio,omitempty" json:"bio"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password"`
	Phone     int64              `bson:"phone,omitempty" json:"phone"`
	PhotoURL  string             `bson:"photo_url,omitempty" json:"photoURL"`
	IsActive  bool               `bson:"is_active,omitempty" json:"isActive"`
	CreatedAt int64              `bson:"created_at,omitempty" json:"createdAt"`
	UpdatedAt int64              `bson:"updated_at,omitempty" json:"updated_at"`
}

type UserDTO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name,omitempty" json:"name"`
	Bio      string             `bson:"bio,omitempty" json:"bio"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Phone    int64              `bson:"phone,omitempty" json:"phone"`
	PhotoURL string             `bson:"photo_url,omitempty" json:"photoURL"`
}

type EmailInfo struct {
	Address string
	Subject string
	Body    string
}

type CustomClaims struct {
	TokenType string
	ID        string `json:"id"`
	*jwt.RegisteredClaims
}
