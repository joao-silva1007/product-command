package domain

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	CreateProduct EventType = "CREATE_PRODUCT"
	UpdateProduct EventType = "UPDATE_PRODUCT"
	DeleteProduct EventType = "DELETE_PRODUCT"
)

type Event struct {
	EventID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Identifier string             `json:"identifier" bson:"identifier"`
	Type       EventType          `json:"type" bson:"type"`
	Body       []byte             `json:"body" bson:"body"`
}

func NewEvent(eventType EventType, body []byte) *Event {
	event := new(Event)

	event.Identifier = uuid.NewString()
	event.Type = eventType
	event.Body = body

	return event
}
