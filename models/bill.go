package models

import (
	"github.com/Dparty/dao/restaurant"
)

type Bill struct {
	entity restaurant.Bill
}

type Specification struct {
	ItemId  string            `json:"itemId"`
	Options []restaurant.Pair `json:"options"`
}
