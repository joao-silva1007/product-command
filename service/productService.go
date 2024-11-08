package service

import (
	"context"
	"encoding/json"
	"pt/isep/insis/product-command/domain"
	"pt/isep/insis/product-command/dto"
	"pt/isep/insis/product-command/messages"
	"pt/isep/insis/product-command/publisher"
	"pt/isep/insis/product-command/repository"
	"pt/isep/insis/product-command/utils"
)

type ProductService interface {
	CreateProduct(ctx context.Context, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error)
	CreateProductAndPublish(ctx context.Context, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error)
	DeleteProductBySku(ctx context.Context, sku string) *utils.Error
	DeleteProductBySkuAndPublish(ctx context.Context, sku string) *utils.Error
	UpdateProductBySku(ctx context.Context, sku string, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error)
	UpdateProductBySkuAndPublish(ctx context.Context, sku string, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error)
}

type productServiceStruct struct {
	ProductRepository repository.ProductRepository
	EventRepository   repository.EventRepository
	ProductPublisher  publisher.ProductPublisher
}

func NewProductService(productRepository repository.ProductRepository, eventRepository repository.EventRepository, productPublisher publisher.ProductPublisher) ProductService {
	return &productServiceStruct{productRepository, eventRepository, productPublisher}
}

func (s *productServiceStruct) CreateProduct(ctx context.Context, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error) {
	jsonBytes, _ := json.Marshal(productDto)
	product, err := domain.NewProduct(jsonBytes)

	if err != nil {
		return nil, err
	}

	createdProduct, err := s.ProductRepository.Create(ctx, product)

	if err != nil {
		return nil, err
	}

	jsonBytes, _ = json.Marshal(createdProduct)

	createdProductDTO, err := dto.NewProductDTO(jsonBytes)

	if err != nil {
		return nil, err
	}

	return createdProductDTO, nil
}

func (s *productServiceStruct) CreateProductAndPublish(ctx context.Context, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error) {
	createdProductDto, err := s.CreateProduct(ctx, productDto)

	if err != nil {
		return nil, err
	}

	bytes, bytesErr := json.Marshal(&messages.CreateProductMessage{Product: productDto})
	if bytesErr != nil {
		return nil, &utils.Error{BaseError: bytesErr, StatusCodeToReturn: 400}
	}

	event := domain.NewEvent(domain.CreateProduct, bytes)
	_, err = s.EventRepository.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	err = s.ProductPublisher.PublishMessage(ctx, event)
	if err != nil {
		return nil, err
	}

	return createdProductDto, nil
}

func (s *productServiceStruct) DeleteProductBySku(ctx context.Context, sku string) *utils.Error {
	return s.ProductRepository.DeleteBySku(ctx, sku)
}

func (s *productServiceStruct) DeleteProductBySkuAndPublish(ctx context.Context, sku string) *utils.Error {
	err := s.DeleteProductBySku(ctx, sku)
	if err != nil {
		return err
	}

	bytes, bytesErr := json.Marshal(&messages.DeleteProductMessage{SKU: sku})
	if bytesErr != nil {
		return &utils.Error{BaseError: bytesErr, StatusCodeToReturn: 400}
	}

	event := domain.NewEvent(domain.DeleteProduct, bytes)
	_, err = s.EventRepository.Create(ctx, event)
	if err != nil {
		return err
	}

	return s.ProductPublisher.PublishMessage(ctx, event)
}

func (s *productServiceStruct) UpdateProductBySku(ctx context.Context, sku string, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error) {
	jsonBytes, _ := json.Marshal(productDto)
	product, err := domain.NewProduct(jsonBytes)

	if err != nil {
		return nil, err
	}

	updatedProduct, err := s.ProductRepository.UpdateBySku(ctx, product, sku)

	if err != nil {
		return nil, err
	}

	jsonBytes, _ = json.Marshal(updatedProduct)

	updatedProductDto, err := dto.NewProductDTO(jsonBytes)

	if err != nil {
		return nil, err
	}

	return updatedProductDto, nil
}

func (s *productServiceStruct) UpdateProductBySkuAndPublish(ctx context.Context, sku string, productDto *dto.ProductDTO) (*dto.ProductDTO, *utils.Error) {
	createdProductDto, err := s.UpdateProductBySku(ctx, sku, productDto)

	if err != nil {
		return nil, err
	}

	bytes, bytesErr := json.Marshal(&messages.UpdateProductMessage{Product: createdProductDto, SKU: sku})
	if bytesErr != nil {
		return nil, &utils.Error{BaseError: bytesErr, StatusCodeToReturn: 400}
	}

	event := domain.NewEvent(domain.UpdateProduct, bytes)
	_, err = s.EventRepository.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	err = s.ProductPublisher.PublishMessage(ctx, event)

	if err != nil {
		return nil, err
	}

	return createdProductDto, nil
}
