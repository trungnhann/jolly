package domain

import "time"

const (
	TopicPaymentAuthorized = "payment.authorized"
	TopicPaymentFailed     = "payment.failed"
)

type PaymentAuthorizedEvent struct {
	OrderID   string    `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"occurred_at"`
}

func (e PaymentAuthorizedEvent) EventName() string {
	return TopicPaymentAuthorized
}

func (e PaymentAuthorizedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

type PaymentFailedEvent struct {
	OrderID   string    `json:"order_id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"occurred_at"`
}

func (e PaymentFailedEvent) EventName() string {
	return TopicPaymentFailed
}

func (e PaymentFailedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}
