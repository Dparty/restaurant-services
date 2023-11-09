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

func (b Bill) Entity() restaurant.Bill {
	return b.entity
}

func (b Bill) Orders() restaurant.Orders {
	return b.entity.Orders
}

func (b Bill) PickUpCode() int64 {
	return b.entity.PickUpCode
}

func (b Bill) Finish(offset int64) {
	b.entity.Status = "FINISH"
	b.entity.Offset = offset
	billRepository.Save(&b.entity)
}
func (b Bill) Set(status string, offset int64) {
	b.entity.Status = status
	b.entity.Offset = offset
	billRepository.Save(&b.entity)
}

func (b *Bill) Save() {
	billRepository.Save(&b.entity)
}

func (b Bill) OwnerId() uint {
	restaurant := restaurantRepository.GetById(b.entity.RestaurantId)
	return restaurant.Owner().ID()
}

type Specification struct {
	ItemId  string            `json:"itemId"`
	Options []restaurant.Pair `json:"options"`
}
