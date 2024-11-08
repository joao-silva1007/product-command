package domain

import (
	"encoding/json"
	"pt/isep/insis/product-command/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SKU         string             `json:"sku" bson:"sku"`
	Designation string             `json:"designation" bson:"designation"`
	Description string             `json:"description" bson:"description"`
}

func NewProduct(jsonBody []byte) (*Product, *utils.Error) {
	product := new(Product)
	err := json.Unmarshal(jsonBody, &product)

	if err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	return product, nil
}
