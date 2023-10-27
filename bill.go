package restaurantservice

import (
	"time"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/restaurant-services/models"
	"gorm.io/gorm"
)

func NewBillService(inject *gorm.DB) BillService {
	return BillService{restaurantDao.NewBillRepository(inject)}
}

type BillService struct {
	billRepository restaurantDao.BillRepository
}

func PairsToMap(s []restaurantDao.Pair) map[string]string {
	output := make(map[string]string)
	for _, option := range s {
		output[option.Left] = option.Right
	}
	return output
}

func (b BillService) CreateBill(table models.Table, specifications []models.Specification) (*models.Bill, error) {
	var orders restaurantDao.Orders
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
	res := restaurantRepository.GetById(table.Owner().ID())
	entity := restaurantDao.Bill{
		RestaurantId: table.Owner().ID(),
		TableId:      table.ID(),
		Status:       "SUBMITTED",
		Orders:       orders,
		PickUpCode:   pickUpCode,
		TableLabel:   table.Label(),
		Offset:       0,
	}
	b.billRepository.Save(&entity)
	bill := models.NewBill(entity)
	models.PrintBill(res.Printers(), res.Name, bill.Entity(), table.Entity(), false)
	return &bill, nil
}

func (b BillService) ListBills(restaurantId uint, tableId *uint, status *string, startAt, endAt *time.Time) []models.Bill {
	ctx := db.Model(&restaurantDao.Bill{})
	ctx = ctx.Where("restaurant_id = ?", restaurantId)
	if tableId != nil {
		ctx = ctx.Where("table_id = ?", *tableId)
	}
	if status != nil {
		ctx = ctx.Where("status = ?", *status)
	}
	if startAt != nil {
		ctx = ctx.Where("created_at >= ?", *startAt)
	}
	if endAt != nil {
		ctx = ctx.Where("created_at <= ?", *endAt)
	}
	var bs []restaurantDao.Bill
	ctx.Find(&bs)
	var bills []models.Bill
	for _, b := range bs {
		bills = append(bills, models.NewBill(b))
	}
	return bills
}
