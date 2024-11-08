package dto

import (
	"encoding/json"
	"pt/isep/insis/product-command/utils"
)

type ProductDTO struct {
	SKU         string `json:"sku" bson:"sku"`
	Designation string `json:"designation" bson:"designation"`
	Description string `json:"description" bson:"description"`
}

func NewProductDTO(jsonBody []byte) (*ProductDTO, *utils.Error) {
	product := new(ProductDTO)
	err := json.Unmarshal(jsonBody, &product)

	if err != nil {
		return nil, &utils.Error{BaseError: err, StatusCodeToReturn: 400}
	}

	return product, nil
}
