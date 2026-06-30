package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"jolly/backend/common"
	"jolly/backend/common/file"
	"jolly/backend/products/app/command"
	"jolly/backend/products/app/query"
	"jolly/backend/products/domain"
)

type Handler struct {
	commands *command.Handlers
	queries  *query.Handlers
	storage  file.Storage
}

func NewHandler(commands *command.Handlers, queries *query.Handlers, storage file.Storage) *Handler {
	if commands == nil {
		panic("products command handlers cannot be nil")
	}
	if queries == nil {
		panic("products query handlers cannot be nil")
	}
	if storage == nil {
		panic("storage cannot be nil")
	}

	return &Handler{
		commands: commands,
		queries:  queries,
		storage:  storage,
	}
}

func mapUUIDPointer(u *openapi_types.UUID) *common.UUID {
	if u == nil {
		return nil
	}
	res := common.UUID(*u)
	return &res
}

func (h Handler) CreateProduct(ctx context.Context, request CreateProductRequestObject) (CreateProductResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	var variants []command.CreateProductVariant
	if request.Body.Variants != nil {
		variants = make([]command.CreateProductVariant, 0, len(*request.Body.Variants))
		for _, v := range *request.Body.Variants {
			variants = append(variants, command.CreateProductVariant{
				SKU:        v.Sku,
				Name:       v.Name,
				PriceCents: v.PriceCents,
			})
		}
	}

	p, err := h.commands.CreateProduct(ctx, command.CreateProduct{
		Name:         request.Body.Name,
		Description:  request.Body.Description,
		Status:       request.Body.Status,
		CategoryUUID: request.Body.CategoryUuid,
		BrandUUID:    request.Body.BrandUuid,
		Variants:     variants,
	})
	if err != nil {
		return nil, err
	}

	return CreateProduct201JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) ListProducts(ctx context.Context, request ListProductsRequestObject) (ListProductsResponseObject, error) {
	products, err := h.queries.ListProducts(ctx, query.ListProducts{})
	if err != nil {
		return nil, err
	}

	resp := make([]Product, 0, len(products))
	for _, p := range products {
		resp = append(resp, mapProductResponse(p))
	}

	return ListProducts200JSONResponse(resp), nil
}

func (h Handler) GetProduct(ctx context.Context, request GetProductRequestObject) (GetProductResponseObject, error) {
	p, err := h.queries.GetProduct(ctx, query.GetProduct{ID: request.ProductUuid})
	if err != nil {
		return nil, err
	}

	return GetProduct200JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) UpdateProduct(ctx context.Context, request UpdateProductRequestObject) (UpdateProductResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	p, err := h.commands.UpdateProduct(ctx, command.UpdateProduct{
		ID:           request.ProductUuid,
		Name:         request.Body.Name,
		Description:  request.Body.Description,
		Status:       request.Body.Status,
		CategoryUUID: request.Body.CategoryUuid,
		BrandUUID:    request.Body.BrandUuid,
	})
	if err != nil {
		return nil, err
	}

	return UpdateProduct200JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) DeleteProduct(ctx context.Context, request DeleteProductRequestObject) (DeleteProductResponseObject, error) {
	err := h.commands.DeleteProduct(ctx, command.DeleteProduct{ID: request.ProductUuid})
	if err != nil {
		return nil, err
	}

	return DeleteProduct204Response{}, nil
}

func (h Handler) AddVariant(ctx context.Context, request AddVariantRequestObject) (AddVariantResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	p, err := h.commands.AddVariant(ctx, command.AddVariant{
		ProductID:  request.ProductUuid,
		SKU:        request.Body.Sku,
		Name:       request.Body.Name,
		PriceCents: request.Body.PriceCents,
	})
	if err != nil {
		return nil, err
	}

	return AddVariant201JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) UpdateVariant(ctx context.Context, request UpdateVariantRequestObject) (UpdateVariantResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	p, err := h.commands.UpdateVariant(ctx, command.UpdateVariant{
		ProductID:   request.ProductUuid,
		VariantUUID: request.VariantUuid,
		SKU:         request.Body.Sku,
		Name:        request.Body.Name,
		PriceCents:  request.Body.PriceCents,
	})
	if err != nil {
		return nil, err
	}

	return UpdateVariant200JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) DeleteVariant(ctx context.Context, request DeleteVariantRequestObject) (DeleteVariantResponseObject, error) {
	p, err := h.commands.RemoveVariant(ctx, command.RemoveVariant{
		ProductID:   request.ProductUuid,
		VariantUUID: request.VariantUuid,
	})
	if err != nil {
		return nil, err
	}

	return DeleteVariant200JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) UploadVariantImage(ctx context.Context, request UploadVariantImageRequestObject) (UploadVariantImageResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	content, fileName, err := parseMultipartFile(request.Body)
	if err != nil {
		return nil, common.NewInvalidInputError("invalid_file", "%s", err.Error())
	}

	ext := filepath.Ext(fileName)
	imageName := fmt.Sprintf("%s%s", common.NewUUIDv7(), ext)

	url, err := h.storage.StoreFile(ctx, fmt.Sprintf("products/%s", imageName), content)
	if err != nil {
		return nil, err
	}

	position := 0
	if request.Params.Position != nil {
		position = *request.Params.Position
	}

	p, err := h.commands.AddVariantImage(ctx, command.AddVariantImage{
		ProductID: request.ProductUuid,
		VariantID: request.VariantUuid,
		URL:       url,
		Position:  position,
	})
	if err != nil {
		return nil, err
	}

	return UploadVariantImage201JSONResponse(mapProductResponse(p)), nil
}

func (h Handler) DeleteVariantImage(ctx context.Context, request DeleteVariantImageRequestObject) (DeleteVariantImageResponseObject, error) {
	p, err := h.commands.RemoveVariantImage(ctx, command.RemoveVariantImage{
		ProductID: request.ProductUuid,
		VariantID: request.VariantUuid,
		ImageID:   request.ImageUuid,
	})
	if err != nil {
		return nil, err
	}

	return DeleteVariantImage200JSONResponse(mapProductResponse(p)), nil
}

// Categories

func (h Handler) CreateCategory(ctx context.Context, request CreateCategoryRequestObject) (CreateCategoryResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	c, err := h.commands.CreateCategory(ctx, command.CreateCategory{
		ParentCategoryUUID: request.Body.ParentCategoryUuid,
		Name:               request.Body.Name,
		Slug:               request.Body.Slug,
	})
	if err != nil {
		return nil, err
	}

	return CreateCategory201JSONResponse(mapCategoryResponse(c)), nil
}

func (h Handler) ListCategories(ctx context.Context, request ListCategoriesRequestObject) (ListCategoriesResponseObject, error) {
	cats, err := h.queries.ListCategories(ctx, query.ListCategories{})
	if err != nil {
		return nil, err
	}

	resp := make([]Category, 0, len(cats))
	for _, c := range cats {
		resp = append(resp, mapCategoryResponse(c))
	}

	return ListCategories200JSONResponse(resp), nil
}

// Brands

func (h Handler) CreateBrand(ctx context.Context, request CreateBrandRequestObject) (CreateBrandResponseObject, error) {
	if request.Body == nil {
		return nil, common.NewInvalidInputError("empty_body", "request body is required")
	}

	b, err := h.commands.CreateBrand(ctx, command.CreateBrand{
		Name: request.Body.Name,
		Slug: request.Body.Slug,
	})
	if err != nil {
		return nil, err
	}

	return CreateBrand201JSONResponse(mapBrandResponse(b)), nil
}

func (h Handler) ListBrands(ctx context.Context, request ListBrandsRequestObject) (ListBrandsResponseObject, error) {
	brands, err := h.queries.ListBrands(ctx, query.ListBrands{})
	if err != nil {
		return nil, err
	}

	resp := make([]Brand, 0, len(brands))
	for _, b := range brands {
		resp = append(resp, mapBrandResponse(b))
	}

	return ListBrands200JSONResponse(resp), nil
}

// Mappers

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers, storage file.Storage) error {
	handler := NewHandler(commands, queries, storage)
	RegisterHandlers(e, NewStrictHandler(handler, nil))
	return nil
}

func mapProductResponse(p domain.Product) Product {
	variants := make([]Variant, 0, len(p.Variants()))
	for _, v := range p.Variants() {
		images := make([]VariantImage, 0, len(v.Images()))
		for _, img := range v.Images() {
			images = append(images, VariantImage{
				ImageUuid:   img.ID(),
				VariantUuid: img.VariantID(),
				Url:         img.URL(),
				Position:    int(img.Position()),
				CreatedAt:   img.CreatedAt(),
			})
		}

		variants = append(variants, Variant{
			VariantUuid: v.ID(),
			ProductUuid: v.ProductID(),
			Sku:         v.SKU(),
			Name:        v.Name(),
			PriceCents:  v.PriceCents(),
			Images:      images,
			CreatedAt:   v.CreatedAt(),
			UpdatedAt:   v.UpdatedAt(),
		})
	}

	return Product{
		ProductUuid:  p.ID(),
		Name:         p.Name(),
		Description:  p.Description(),
		Status:       p.Status(),
		CategoryUuid: p.CategoryUUID(),
		BrandUuid:    p.BrandUUID(),
		Variants:     variants,
		CreatedAt:    p.CreatedAt(),
		UpdatedAt:    p.UpdatedAt(),
	}
}

func mapCategoryResponse(c domain.Category) Category {
	return Category{
		CategoryUuid:       c.ID(),
		ParentCategoryUuid: c.ParentCategoryUUID(),
		Name:               c.Name(),
		Slug:               c.Slug(),
		CreatedAt:          c.CreatedAt(),
		UpdatedAt:          c.UpdatedAt(),
	}
}

func mapBrandResponse(b domain.Brand) Brand {
	return Brand{
		BrandUuid: b.ID(),
		Name:      b.Name(),
		Slug:      b.Slug(),
		CreatedAt: b.CreatedAt(),
		UpdatedAt: b.UpdatedAt(),
	}
}

func parseMultipartFile(reader *multipart.Reader) ([]byte, string, error) {
	for {
		part, err := reader.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, "", err
		}

		if part.FormName() == "file" {
			defer part.Close()
			content, err := io.ReadAll(part)
			if err != nil {
				return nil, "", err
			}
			return content, part.FileName(), nil
		}
	}
	return nil, "", errors.New("no file part named 'file' found")
}
