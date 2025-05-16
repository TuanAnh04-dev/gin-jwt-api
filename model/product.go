package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `Json:"name", omitempty, bson:"name"`
	Price    float64            `Json:"price", omitempty, bson:"price"`
	Quantity int                `Json:"quantity", omitempty, bson:"quantity"`
}
