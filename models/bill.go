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

type Specification struct {
	ItemId  string            `json:"itemId"`
	Options []restaurant.Pair `json:"options"`
}
