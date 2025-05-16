package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `Json:"id,omitempty" bson:"_id,omitempty"` // tag golang
	Email       string             `Json:"email" bson:"email"`
	Password    string             `Json:"password" bson:"password"`
	DisplayName string             `Json:"displayName" bson:"displayName"`
}
