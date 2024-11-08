package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"pt/isep/insis/product-command/domain"
	"pt/isep/insis/product-command/utils"
)

type EventRepository interface {
	findByIdentifier(ctx context.Context, identifier string) (*domain.Event, *utils.Error)
	Create(ctx context.Context, event *domain.Event) (*domain.Event, *utils.Error)
}

type EventRepositoryStruct struct {
	coll *mongo.Collection
}

func NewEventRepository(coll *mongo.Collection) *EventRepositoryStruct {
	return &EventRepositoryStruct{coll}
}

func (repo *EventRepositoryStruct) findByIdentifier(ctx context.Context, identifier string) (*domain.Event, *utils.Error) {
	filter := bson.M{"identifier": identifier}
	event := &domain.Event{}
	err := repo.coll.FindOne(ctx, filter).Decode(event)

	if err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 404}
	}

	return event, nil
}

func (repo *EventRepositoryStruct) Create(ctx context.Context, event *domain.Event) (*domain.Event, *utils.Error) {
	if _, err := repo.coll.InsertOne(ctx, event); err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	} else {
		return repo.findByIdentifier(ctx, event.Identifier)
	}
}
