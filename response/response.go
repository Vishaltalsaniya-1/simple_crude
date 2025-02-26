package response

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StdResponse struct {
	ID          primitive.ObjectID `json:"id"`
	AuthID      primitive.ObjectID `json:"auth_id" bson:"auth_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Tag         []string           `json:"tag" bson:"tag"`
	Student     bool               `json:"student" bson:"student"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at,omitempty"`
}
