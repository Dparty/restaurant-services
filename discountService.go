package restaurantservice

import (
	restaurantDao "github.com/Dparty/dao/restaurant"
)

var discountService *DiscountService

func GetDiscountService() *DiscountService {
	if discountService == nil {
		discountService = NewDiscountService()
	}
	return discountService
}

type DiscountService struct {
	discountRepository *restaurantDao.DiscountRepository
}

func NewDiscountService() *DiscountService {
	return &DiscountService{restaurantDao.GetDiscountRepository()}
}

// Delete Discount
func (d DiscountService) DeleteDiscount(id uint) {
	discount := d.discountRepository.Find(id)
	if discount != nil {
		d := NewDiscount(*discount)
		d.Delete()
	}
}

func (d DiscountService) GetDiscount(id uint) *Discount {
	discount := d.discountRepository.Find(id)
	if discount != nil {
		return &Discount{entity: *discount}
	}
	return nil
}
