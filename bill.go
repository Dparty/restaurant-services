package restaurantservice

import (
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/restaurant-services/models"
	"gorm.io/gorm"
)

func NewBillService(inject *gorm.DB) BillService {
	return BillService{restaurant.NewBillRepository(inject)}
}

type BillService struct {
	billRepository restaurant.BillRepository
}

func PairsToMap(s []restaurant.Pair) map[string]string {
	output := make(map[string]string)
	for _, option := range s {
		output[option.Left] = option.Right
	}
	return output
}

func (b BillService) CreateBill(table models.Table, specifications []models.Specification) (*models.Bill, error) {
	var orders restaurant.Orders
	for _, specification := range specifications {
		item := itemRepository.GetById(utils.StringToUint(specification.ItemId))
		if item == nil {
			return nil, fault.ErrNotFound
		}
		order, err := item.CreateOrder(specification.Options)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	pickUpCode := restaurantRepository.GetById(table.Owner().ID()).PickUpCode()
	entity := restaurant.Bill{
		RestaurantId: table.Owner().ID(),
		TableId:      table.ID(),
		Status:       "SUBMIT",
		Orders:       orders,
		PickUpCode:   pickUpCode,
		TableLabel:   table.Label(),
	}
	b.billRepository.Save(&entity)
	bill := models.NewBill(entity)
	return &bill, nil
}
