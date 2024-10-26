package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Age       int                `bson:"age,omitempty" json:"age,omitempty"`
}

// UpdateUserRequest is used for updating user data
type UpdateUserRequest struct {
	Filter map[string]interface{} `json:"filter"`
	Update map[string]interface{} `json:"update"`
}

// FilterRequest is used for deleting user data
type FilterRequest struct {
	Filter map[string]interface{} `json:"filter"` // e.g., {"firstName": "John"}
}
