package models

import (
	"github.com/Dparty/dao/restaurant"
)

func NewBill(entity restaurant.Bill) Bill {
	return Bill{entity: entity}
}

type Bill struct {
	entity restaurant.Bill
}

func (b Bill) ID() uint {
	return b.entity.ID
}

func (b Bill) Orders() restaurant.Orders {
	return b.entity.Orders
}

func (b Bill) PickUpCode() int64 {
	return b.entity.PickUpCode
}

type Specification struct {
	ItemId  string            `json:"itemId"`
	Options []restaurant.Pair `json:"options"`
}
