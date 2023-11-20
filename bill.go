package restaurantservice

import (
	"github.com/Dparty/common/utils"
	restaurantDao "github.com/Dparty/dao/restaurant"
)

func NewBill(entity restaurantDao.Bill) Bill {
	return Bill{entity: entity}
}

type Bill struct {
	entity restaurantDao.Bill
}

func (b Bill) ID() uint {
	return b.entity.ID
}

func (b Bill) Entity() restaurantDao.Bill {
	return b.entity
}

func (b Bill) Orders() restaurantDao.Orders {
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

func (b Bill) CancelItems(specifications []Specification) {
	var newOrders restaurantDao.Orders
	copy(newOrders, b.entity.Orders)
	// for _, specification := range specifications {
	// 	for _, order := range b.entity.Orders {

	// 	}
	// }
	// b.entity.Orders
}

type Specification struct {
	ItemId  string               `json:"itemId"`
	Options []restaurantDao.Pair `json:"options"`
}

func (s Specification) Equal(order restaurantDao.Order) bool {
	if utils.StringToUint(s.ItemId) != order.Item.ID() {
		return false
	}
	for _, o := range s.Options {
		for _, o2 := range order.Specification {
			if o.Left == o2.Left && o.Right != o2.Right {
				return false
			}
		}
	}
	return true
}
