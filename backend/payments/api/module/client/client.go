package client

import "context"

type Payments interface {
	Authorize(ctx context.Context, req AuthorizePaymentRequest) (AuthorizePaymentResponse, error)
}

type AuthorizePaymentRequest struct {
	OrderID     string
	CustomerID  string
	AmountCents int64
	Currency    string
}

type AuthorizePaymentResponse struct {
	PaymentID string
	Status    string
}
