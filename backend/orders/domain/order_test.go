package domain_test

import (
	"testing"
	"time"

	"jolly/backend/common"
	"jolly/backend/common/shared"
	"jolly/backend/orders/domain"
)

func TestOrder_StatusTransitions(t *testing.T) {
	orderID := domain.OrderID{UUID: common.NewUUIDv7()}
	customerID := common.NewUUIDv7().String()
	currency := shared.MustNewCurrency("USD")

	items := []domain.LineItem{
		{
			UUID:           domain.LineItemUUID{UUID: common.NewUUIDv7()},
			SKU:            "SKU-123",
			Quantity:       2,
			UnitPriceCents: 1000,
		},
	}

	t.Run("successful flow: pending_payment -> paid -> confirmed", func(t *testing.T) {
		order, err := domain.NewOrder(orderID, customerID, items, currency, time.Now())
		if err != nil {
			t.Fatalf("failed to construct order: %v", err)
		}

		if order.Status != domain.StatusPendingPayment {
			t.Errorf("expected initial status %s, got %s", domain.StatusPendingPayment, order.Status)
		}

		// Mark paid
		err = order.MarkPaid()
		if err != nil {
			t.Fatalf("failed to mark order paid: %v", err)
		}
		if order.Status != domain.StatusPaid {
			t.Errorf("expected status %s, got %s", domain.StatusPaid, order.Status)
		}

		// Mark confirmed
		err = order.MarkConfirmed()
		if err != nil {
			t.Fatalf("failed to mark order confirmed: %v", err)
		}
		if order.Status != domain.StatusConfirmed {
			t.Errorf("expected status %s, got %s", domain.StatusConfirmed, order.Status)
		}
	})

	t.Run("COD flow: pending_payment -> confirmed -> paid", func(t *testing.T) {
		order, err := domain.NewOrder(orderID, customerID, items, currency, time.Now())
		if err != nil {
			t.Fatalf("failed to construct order: %v", err)
		}

		// Mark confirmed (inventory reserved)
		err = order.MarkConfirmed()
		if err != nil {
			t.Fatalf("failed to mark order confirmed: %v", err)
		}
		if order.Status != domain.StatusConfirmed {
			t.Errorf("expected status %s, got %s", domain.StatusConfirmed, order.Status)
		}

		// Mark paid (payment authorized)
		err = order.MarkPaid()
		if err != nil {
			t.Fatalf("failed to mark order paid: %v", err)
		}
		if order.Status != domain.StatusPaid {
			t.Errorf("expected status %s, got %s", domain.StatusPaid, order.Status)
		}
	})

	t.Run("transition to failed from pending_payment", func(t *testing.T) {
		order, err := domain.NewOrder(orderID, customerID, items, currency, time.Now())
		if err != nil {
			t.Fatalf("failed to construct order: %v", err)
		}

		err = order.MarkFailed()
		if err != nil {
			t.Fatalf("failed to mark failed: %v", err)
		}
		if order.Status != domain.StatusFailed {
			t.Errorf("expected status %s, got %s", domain.StatusFailed, order.Status)
		}
	})

	t.Run("transition to failed from paid", func(t *testing.T) {
		order, err := domain.NewOrder(orderID, customerID, items, currency, time.Now())
		if err != nil {
			t.Fatalf("failed to construct order: %v", err)
		}

		_ = order.MarkPaid()

		err = order.MarkFailed()
		if err != nil {
			t.Fatalf("failed to mark failed: %v", err)
		}
		if order.Status != domain.StatusFailed {
			t.Errorf("expected status %s, got %s", domain.StatusFailed, order.Status)
		}
	})
}
