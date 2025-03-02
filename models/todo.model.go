package todo_model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TODO struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BODY      string             `json:"body" bson:"body"`
	Completed bool               `json:"completed" bson:"completed"`
}
