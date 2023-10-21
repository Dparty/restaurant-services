package restaurantservice

import (
	"github.com/Dparty/dao/restaurant"
	"gorm.io/gorm"
)

func NewBillService(inject *gorm.DB) BillService {
	return BillService{restaurant.NewBillRepository(inject)}
}

type BillService struct {
	billRepository restaurant.BillRepository
}

func (b BillService) CreateBill() {

}
