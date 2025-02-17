package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Std struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthID      primitive.ObjectID `json:"auth_id" bson:"auth_id" binding:"required"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Description string             `json:"description" bson:"description"`
	Tag         []string           `json:"tag"  bson:"tag"`
	Student     bool               `json:"student" bson:"student"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
