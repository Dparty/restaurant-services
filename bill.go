package restaurantservice

import (
	"github.com/Dparty/common/data"
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

func (b *Bill) Finish(offset int64) {
	b.entity.Status = "FINISH"
	b.entity.Offset = offset
	b.Submit()
}
func (b *Bill) Set(status string, offset int64) {
	b.entity.Status = status
	b.entity.Offset = offset
	b.Submit()
}

func (b *Bill) Submit() {
	billRepository.Save(&b.entity)
}

func (b Bill) OwnerId() uint {
	restaurant := restaurantRepository.GetById(b.entity.RestaurantId)
	return restaurant.Owner().ID()
}

func (b *Bill) CancelItems(specifications []Specification) {
	if len(specifications) == 0 {
		return
	}
	head := specifications[0]
	for i, order := range b.entity.Orders {
		if head.Equal(order) {
			b.entity.Orders = append(b.entity.Orders[:i], b.entity.Orders[i+1:]...)
			break
		}
	}
	b.CancelItems(specifications[1:])
}

type Specification struct {
	ItemId  string                      `json:"itemId"`
	Options []data.Pair[string, string] `json:"options"`
}

func (s Specification) Equal(order restaurantDao.Order) bool {
	if utils.StringToUint(s.ItemId) != order.Item.ID() {
		return false
	}
	for _, o := range s.Options {
		for _, o2 := range order.Specification {
			if o.L == o2.L && o.R != o2.R {
				return false
			}
		}
	}
	return true
}
