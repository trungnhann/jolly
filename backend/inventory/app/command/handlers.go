package command

import (
	"context"
	"fmt"
	"strings"

	"jolly/backend/common"
	"jolly/backend/common/module/contracts"
	"jolly/backend/inventory/domain"
	productsclient "jolly/backend/products/api/module/client"
)

type Handlers struct {
	repo      domain.Repository
	contracts *contracts.Contracts
}

func NewHandlers(repo domain.Repository, contracts *contracts.Contracts) *Handlers {
	if repo == nil {
		panic("inventory repository cannot be nil")
	}
	return &Handlers{
		repo:      repo,
		contracts: contracts,
	}
}

type Reserve struct {
	OrderID string
	Items   []struct {
		SKU      string
		Quantity int
	}
}

func (h *Handlers) Reserve(ctx context.Context, cmd Reserve) error {
	if strings.TrimSpace(cmd.OrderID) == "" {
		return fmt.Errorf("order id cannot be empty")
	}
	if len(cmd.Items) == 0 {
		return fmt.Errorf("at least one stock item is required")
	}

	var orderUUID common.UUID
	if err := orderUUID.UnmarshalText([]byte(cmd.OrderID)); err != nil {
		return fmt.Errorf("invalid order id: %w", err)
	}

	items := make([]struct {
		SKU      string
		Quantity int
	}, 0, len(cmd.Items))

	for _, item := range cmd.Items {
		if strings.TrimSpace(item.SKU) == "" {
			return fmt.Errorf("stock item sku cannot be empty")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("stock item quantity must be positive")
		}
		items = append(items, struct {
			SKU      string
			Quantity int
		}{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	reservationUUID := common.NewUUIDv7()
	return h.repo.Reserve(ctx, reservationUUID, orderUUID, items)
}

type UpsertStock struct {
	SKU      string
	Quantity int
}

func (h *Handlers) UpsertStock(ctx context.Context, cmd UpsertStock) (domain.Stock, error) {
	if strings.TrimSpace(cmd.SKU) == "" {
		return domain.Stock{}, fmt.Errorf("sku cannot be empty")
	}
	if cmd.Quantity < 0 {
		return domain.Stock{}, fmt.Errorf("quantity cannot be negative")
	}

	// Verify SKU exists in products module catalog
	if h.contracts != nil && h.contracts.Products != nil {
		_, err := h.contracts.Products.GetVariantBySKU(ctx, productsclient.GetVariantBySKURequest{
			SKU: cmd.SKU,
		})
		if err != nil {
			return domain.Stock{}, common.NewInvalidInputError("sku_not_found", "sku %s does not exist in products catalog", cmd.SKU)
		}
	}

	stockUUID := common.NewUUIDv7()
	return h.repo.UpsertStock(ctx, stockUUID, cmd.SKU, cmd.Quantity)
}
