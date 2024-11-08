package messages

import (
	"pt/isep/insis/product-command/dto"
)

type ProductRpcMessage struct {
	Products []*dto.ProductDTO `json:"products" bson:"products"`
}
