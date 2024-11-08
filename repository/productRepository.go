package repository

import (
	"context"
	"errors"
	"pt/isep/insis/product-command/domain"
	"pt/isep/insis/product-command/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type ProductRepository interface {
	findBySku(ctx context.Context, sku string) (*domain.Product, *utils.Error)
	Create(ctx context.Context, product *domain.Product) (*domain.Product, *utils.Error)
	DeleteBySku(ctx context.Context, sku string) *utils.Error
	UpdateBySku(ctx context.Context, product *domain.Product, sku string) (*domain.Product, *utils.Error)
}

type ProductRepositoryStruct struct {
	coll *mongo.Collection
}

func NewProductRepository(coll *mongo.Collection) *ProductRepositoryStruct {
	return &ProductRepositoryStruct{coll}
}

func (repo *ProductRepositoryStruct) findBySku(ctx context.Context, sku string) (*domain.Product, *utils.Error) {
	filter := bson.M{"sku": sku}
	product := &domain.Product{}
	err := repo.coll.FindOne(ctx, filter).Decode(product)

	if err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 404}
	}

	return product, nil
}

func (repo *ProductRepositoryStruct) Create(ctx context.Context, product *domain.Product) (*domain.Product, *utils.Error) {
	prod, _ := repo.findBySku(ctx, product.SKU)

	if prod != nil {
		return nil, &utils.Error{BaseError: errors.New("product with that sku already exists"), StatusCodeToReturn: 409}
	}

	if _, err := repo.coll.InsertOne(ctx, product); err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	} else {
		return repo.findBySku(ctx, product.SKU)
	}
}

func (repo *ProductRepositoryStruct) DeleteBySku(ctx context.Context, sku string) *utils.Error {
	filter := bson.M{"sku": sku}

	res, err := repo.coll.DeleteOne(ctx, filter)

	if err != nil {
		return &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	if res.DeletedCount == 0 {
		return &utils.Error{BaseError: errors.New("no product deleted"), StatusCodeToReturn: 400}
	}

	return nil
}

func (repo *ProductRepositoryStruct) UpdateBySku(ctx context.Context, product *domain.Product, sku string) (*domain.Product, *utils.Error) {
	res, err := repo.coll.UpdateOne(ctx, bson.M{"sku": sku}, bson.M{"$set": bson.M{"designation": product.Designation, "description": product.Description}})

	if err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	if res.ModifiedCount == 0 {
		return nil, &utils.Error{BaseError: errors.New("no product updated"), StatusCodeToReturn: 400}
	}

	return repo.findBySku(ctx, sku)
}
