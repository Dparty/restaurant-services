package restaurantservice

import (
	"fmt"
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

func (b BillService) SetBill(ownerId uint, billIdList []uint, offset int64, status string) error {
	if len(billIdList) == 0 {
		return nil
	}
	var billsDao []restaurantDao.Bill
	db.Find(&billsDao, billIdList)
	var bills []models.Bill
	for _, bill := range billsDao {
		bills = append(bills, models.NewBill(bill))
	}
	for _, bill := range bills {
		if bill.OwnerId() != ownerId {
			return fault.ErrPermissionDenied
		}
	}
	for _, bill := range bills {
		bill.Set(status, offset)
	}
	return nil
}

func (b BillService) PrintBills(ownerId uint, billIdList []uint, offset int64) error {
	if len(billIdList) == 0 {
		return nil
	}
	var billsDao []restaurantDao.Bill
	db.Find(&billsDao, billIdList)
	var bills []models.Bill
	for _, bill := range billsDao {
		bills = append(bills, models.NewBill(bill))
	}
	for _, bill := range bills {
		if bill.OwnerId() != ownerId {
			return fault.ErrPermissionDenied
		}
	}
	restaurant := restaurantRepository.GetById(bills[0].Entity().RestaurantId)
	fmt.Println(restaurant)
	printers := restaurant.Printers()
	fmt.Println(printers)
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurant.Name)
	// content += models.FinishString(
	// 	offset,
	// 	golambda.Map(bills,
	// 		func(_ int, b models.Bill) restaurantDao.Bill {
	// 			return b.Entity()
	// 		}))
	// for _, printer := range printers {
	// 	if printer.Type == "BILL" {
	// 		p, _ := printerFactory.Connect(printer.Sn)
	// 		p.Print(content, "")
	// 	}
	// }
	return nil
}
