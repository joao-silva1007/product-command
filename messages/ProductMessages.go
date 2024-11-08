package messages

import (
	"pt/isep/insis/product-command/dto"
)

type CreateProductMessage struct {
	Product *dto.ProductDTO `json:"product" bson:"product"`
}

type UpdateProductMessage struct {
	Product *dto.ProductDTO `json:"product" bson:"product"`
	SKU     string          `json:"sku" bson:"sku"`
}

type DeleteProductMessage struct {
	SKU string `json:"sku" bson:"sku"`
}
