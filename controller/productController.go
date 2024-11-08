package controller

import (
	"pt/isep/insis/product-command/dto"
	"pt/isep/insis/product-command/helpers"
	"pt/isep/insis/product-command/service"

	"github.com/gofiber/fiber/v2"
)

type ProductController interface {
	CreateProduct(ctx *fiber.Ctx) error
	DeleteProductBySku(ctx *fiber.Ctx) error
	UpdateProductBySku(ctx *fiber.Ctx) error
}

type ProductControllerStruct struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductControllerStruct {
	return &ProductControllerStruct{productService}
}

func (c *ProductControllerStruct) CreateProduct(ctx *fiber.Ctx) error {
	productDTO, err := dto.NewProductDTO(ctx.Body())

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	createdProductDTO, err := c.ProductService.CreateProductAndPublish(ctx.UserContext(), productDTO)

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	return helpers.ResponseWithBodyAndStatusCode(ctx, createdProductDTO, fiber.StatusCreated)
}

func (c *ProductControllerStruct) DeleteProductBySku(ctx *fiber.Ctx) error {
	sku := ctx.Params("sku")

	err := c.ProductService.DeleteProductBySkuAndPublish(ctx.UserContext(), sku)

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	return helpers.DatabaseResponseWithoutData(ctx)
}

func (c *ProductControllerStruct) UpdateProductBySku(ctx *fiber.Ctx) error {
	sku := ctx.Params("sku")
	productDTO, err := dto.NewProductDTO(ctx.Body())

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	updatedProductDTO, err := c.ProductService.UpdateProductBySkuAndPublish(ctx.UserContext(), sku, productDTO)

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	return helpers.ResponseWithBodyAndStatusCode(ctx, updatedProductDTO, fiber.StatusOK)
}
