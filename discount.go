package restaurantservice

import (
	abstract "github.com/Dparty/common/abstract"
	restaurantDao "github.com/Dparty/dao/restaurant"
)

type Discount struct {
	entity restaurantDao.Discount
}

func NewDiscount(entity restaurantDao.Discount) Discount {
	return Discount{entity: entity}
}

func (d Discount) ID() uint {
	return d.entity.ID
}

func (d Discount) Entity() restaurantDao.Discount {
	return d.entity
}

func (d Discount) Owner() abstract.Owner {
	restaurant := restaurantRepository.GetById(d.entity.RestaurantId)
	return NewRestaurant(*restaurant)
}

func (d *Discount) SetLabel(label string) *Discount {
	d.entity.Label = label
	return d
}

func (d *Discount) SetOffset(offset int64) *Discount {
	d.entity.Offset = offset
	return d
}

func (d *Discount) Submit() {
	discountRepository.Save(&d.entity)
}

func (d *Discount) Delete() {
	discountRepository.Delete(&d.entity)
}
