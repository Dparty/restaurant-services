package restaurantservice

import (
	"time"

	"github.com/Dparty/common/data"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/pubsub"
	"github.com/Dparty/common/utils"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/chenyunda218/golambda"
)

var billService *BillService
var pb = pubsub.GetPubSub()

func GetBillService() *BillService {
	if billService == nil {
		billService = NewBillService()
	}
	return billService
}

func NewBillService() *BillService {
	return &BillService{restaurantDao.GetBillRepository()}
}

type BillService struct {
	billRepository *restaurantDao.BillRepository
}

func PairsToMap(s []data.Pair[string, string]) map[string]string {
	output := make(map[string]string)
	for _, option := range s {
		output[option.L] = option.R
	}
	return output
}

func (b BillService) CreateBill(table Table, specifications []Specification, offset int64) (*Bill, error) {
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
	res, _ := restaurantService.GetRestaurant(table.Owner().ID())
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
	bill := NewBill(entity)
	PrintBill(res.Printers(), res.Name(), bill.Entity(), table.Entity(), offset, false)
	bill.Publish()
	return &bill, nil
}

func (b BillService) GetById(id uint) (Bill, error) {
	entity := b.billRepository.GetById(id)
	if entity == nil {
		return Bill{}, fault.ErrNotFound
	}
	return NewBill(*entity), nil
}

func (b BillService) ListBills(restaurantId uint, tableId *uint, status *string, startAt, endAt *time.Time) []Bill {
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
	var bills []Bill
	for _, b := range bs {
		bills = append(bills, NewBill(b))
	}
	return bills
}

func (b BillService) SetBill(ownerId uint, billIdList []uint, offset int64, status string) error {
	if len(billIdList) == 0 {
		return nil
	}
	var billsDao []restaurantDao.Bill
	db.Find(&billsDao, billIdList)
	var bills []Bill
	for _, bill := range billsDao {
		bills = append(bills, NewBill(bill))
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
	billsDao := billRepository.List(billIdList)
	db.Find(&billsDao, billIdList)
	var bills []Bill
	for _, bill := range billsDao {
		bills = append(bills, NewBill(bill))
	}
	for _, bill := range bills {
		if bill.OwnerId() != ownerId {
			return fault.ErrPermissionDenied
		}
	}
	if len(bills) == 0 {
		return nil
	}
	restaurant, _ := restaurantService.GetRestaurant(bills[0].Entity().RestaurantId)
	printers := restaurant.Printers()
	for _, printer := range printers {
		if printer.Type() == "BILL" {
			var pc feieyun.PrintContent
			pc.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: restaurant.Name()}})
			FinishString(
				&pc,
				offset,
				golambda.Map(bills,
					func(_ int, b Bill) restaurantDao.Bill {
						return b.Entity()
					}), printer.Width())
			p, _ := printerFactory.Connect(printer.Sn())
			p.Print(pc.String(), "")
		}
	}
	return nil
}
