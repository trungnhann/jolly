package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

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
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Status:      request.Body.Status,
		Variants:    variants,
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
		ID:          request.ProductUuid,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Status:      request.Body.Status,
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
		return nil, common.NewInvalidInputError("empty_body", "multipart body is required")
	}

	fileContent, filename, err := parseMultipartFile(request.Body)
	if err != nil {
		return nil, common.NewInvalidInputError("invalid-multipart", "failed to parse multipart file: %v", err)
	}

	productUUID := request.ProductUuid
	variantUUID := request.VariantUuid

	ext := filepath.Ext(filename)
	imageUUID := domain.VariantImageUUID{UUID: common.NewUUIDv7()}
	storagePath := fmt.Sprintf("products/%s/variants/%s/%s%s", productUUID.String(), variantUUID.String(), imageUUID.String(), ext)

	url, err := h.storage.StoreFile(ctx, storagePath, fileContent)
	if err != nil {
		return nil, err
	}

	position := 0
	if request.Params.Position != nil {
		position = *request.Params.Position
	}

	p, err := h.commands.AddVariantImage(ctx, command.AddVariantImage{
		ProductID: productUUID,
		VariantID: variantUUID,
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

func Register(ctx context.Context, e common.EchoRouter, commands *command.Handlers, queries *query.Handlers, storage file.Storage) error {
	_ = ctx

	handler := Handler{
		commands: commands,
		queries:  queries,
		storage:  storage,
	}

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
		ProductUuid: p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Status:      p.Status(),
		Variants:    variants,
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
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
