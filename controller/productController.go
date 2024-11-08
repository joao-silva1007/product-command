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

type productControllerStruct struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &productControllerStruct{productService}
}

func (c *productControllerStruct) CreateProduct(ctx *fiber.Ctx) error {
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

func (c *productControllerStruct) DeleteProductBySku(ctx *fiber.Ctx) error {
	sku := ctx.Params("sku")

	err := c.ProductService.DeleteProductBySkuAndPublish(ctx.UserContext(), sku)

	if err != nil {
		return helpers.ErrorResponse(ctx, err.BaseError.Error(), err.StatusCodeToReturn)
	}

	return helpers.DatabaseResponseWithoutData(ctx)
}

func (c *productControllerStruct) UpdateProductBySku(ctx *fiber.Ctx) error {
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
