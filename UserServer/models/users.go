package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     *string            `json:"email"`
	Password  *string            `json:"password"`
	FirstName *string            `json:"firstname"`
	LastName  *float64           `json:"lastname""`
	Age       *int16             `json:"age"`
	City      *string            `json:"city"`
}
