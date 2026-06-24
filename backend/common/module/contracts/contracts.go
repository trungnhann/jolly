package contracts

import (
	"errors"

	inventoryclient "jolly/backend/inventory/api/module/client"
	ordersclient "jolly/backend/orders/api/module/client"
	paymentsclient "jolly/backend/payments/api/module/client"
	usersclient "jolly/backend/users/api/module/client"
)

type Contracts struct {
	Orders    ordersclient.Orders
	Payments  paymentsclient.Payments
	Inventory inventoryclient.Inventory
	Users     usersclient.Users
}

func (c *Contracts) Verify() error {
	var err error

	if c.Orders == nil {
		err = errors.Join(err, errors.New("orders module contract is empty"))
	}
	if c.Payments == nil {
		err = errors.Join(err, errors.New("payments module contract is empty"))
	}
	if c.Inventory == nil {
		err = errors.Join(err, errors.New("inventory module contract is empty"))
	}
	if c.Users == nil {
		err = errors.Join(err, errors.New("users module contract is empty"))
	}

	return err
}
