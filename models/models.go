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
	UpdatedAt int64              `bson:"updated_at,omitempty" json:"updatedAt"`
}

type UserDTO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name,omitempty" json:"name"`
	Bio      string             `bson:"bio,omitempty" json:"bio"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Phone    int64              `bson:"phone,omitempty" json:"phone"`
	PhotoURL string             `bson:"photo_url,omitempty" json:"photoURL"`
}

type GroupChat struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name,omitempty" json:"name"`
	Description  string             `bson:"description,omitempty" json:"description"`
	ImageURL     string             `bson:"image_url,omitempty" json:"imageURL"`
	GroupOwner   primitive.ObjectID `bson:"group_owner,omitempty" json:"groupOwner"`
	Participants []User             `bson:"participants,omitempty" json:"participants"`
	CreatedAt    int64              `bson:"created_at,omitempty" json:"createdAt"`
	UpdatedAt    int64              `bson:"updated_at,omitempty" json:"updatedAt"`
}

type GroupChatDTO struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	Name         string `bson:"name,omitempty" json:"name"`
	Description  string `bson:"description,omitempty" json:"description"`
	ImageURL     string `bson:"image_url,omitempty" json:"imageURL"`
	GroupOwner   string `bson:"group_owner,omitempty" json:"groupOwner"`
	Participants []User `bson:"participants,omitempty" json:"participants"`
}

type EmailInfo struct {
	Address string
	Subject string
	Body    string
}

type ClusterOfMessages struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BelongsToGroup primitive.ObjectID `bson:"belongs_to_group,omitempty" json:"belongsToGroup"`
	Messages       []MsgContent       `bson:"messages,omitempty" json:"messages"`
	CreatedAt      int64              `bson:"created_at,omitempty" json:"createdAt"`
	UpdatedAt      int64              `bson:"updated_at,omitempty" json:"updatedAt"`
	IsClosed       bool               `bson:"is_closed" json:"isClosed"`
}

type ClusterOfMessagesDTO struct {
	BelongsToGroup string          `bson:"belongs_to_group,omitempty" json:"belongsToGroup"`
	Messages       []MsgContentDTO `bson:"messages,omitempty" json:"messages"`
}

type MsgContent struct {
	From        primitive.ObjectID `json:"from"`
	To          primitive.ObjectID `json:"to"`
	TextContent string             `json:"textContent"`
}

type Operation int

const (
	MessageToGroup Operation = iota
	AddToGroup
)

type MsgContentDTO struct {
	OptType     Operation `json:"optType"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	TextContent string    `json:"textContent"`
}

type CustomClaims struct {
	TokenType string
	ID        string `json:"id"`
	*jwt.RegisteredClaims
}

type FileInfo struct {
	Size        int64
	ContentType string
}
