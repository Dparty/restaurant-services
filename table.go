package restaurantservice

import (
	"fmt"
	"math"
	"time"

	abstract "github.com/Dparty/common/abstract"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/chenyunda218/golambda"
)

type Table struct {
	entity restaurantDao.Table
}

func (t Table) Owner() abstract.Owner {
	return t.entity.Owner()
}

func (t Table) ID() uint {
	return t.entity.ID()
}

func (t Table) Label() string {
	return t.entity.Label
}

func (t Table) X() int64 {
	return t.entity.X
}

func (t Table) Y() int64 {
	return t.entity.Y
}

func (t Table) Entity() restaurantDao.Table {
	return t.entity
}

func (t Table) Bills(status *string) []Bill {
	var bills []restaurantDao.Bill
	db.Where("table_id = ?", t.entity.ID()).Where("status = ?", *status).Find(&bills)
	return golambda.Map(bills, func(_ int, bill restaurantDao.Bill) Bill {
		return NewBill(bill)
	})
}

func (t Table) Delete() bool {
	return tableRepository.Delete(&t.entity).RowsAffected != 0
}

func (t Table) PrintBills(offset int64) {
	restaurant, _ := restaurantService.GetRestaurant(t.Owner().ID())
	printers := restaurant.Printers()
	status := "SUBMITTED"
	bills := t.Bills(&status)
	if len(bills) == 0 {
		return
	}

	for _, printer := range printers {
		if printer.Type() == "BILL" {
			var pc feieyun.PrintContent
			pc.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: restaurant.Name()}})
			pc.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("桌號: %s", t.Label())}})
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
}

func (t Table) Finish(offset int64) {
	status := "SUBMITTED"
	bills := t.Bills(&status)
	for _, bill := range bills {
		bill.Finish(offset)
	}
}

func (t Table) Update(label string, x, y int64) bool {
	tables := tableRepository.List("restaurant_id = ?", t.Owner().ID())
	if len(golambda.Filter(tables, func(_ int, i restaurantDao.Table) bool {
		return i.ID() != t.entity.ID() && (i.Label == label || (i.X == x && i.Y == y))
	})) != 0 {
		return false
	}
	entity := t.Entity()
	entity.Label = label
	entity.X = x
	entity.Y = y
	tableRepository.Save(&entity)
	return true
}

func FinishString(pc *feieyun.PrintContent, offset int64, bills []restaurantDao.Bill, width int) string {
	_offset := float64(offset+100) / 100

	total := 0
	for _, bill := range bills {
		total += int(bill.Total())
		orderNumbers := make([]OrderNumber, 0)
		for _, order := range bill.Orders {
			orderNumbers = MargeOrderHelper(order, orderNumbers)
		}
		pc.AddDiv(width)
		pc.AddLines(&feieyun.Text{Content: fmt.Sprintf("餐號: %d", bill.PickUpCode)})
		pc.AddLines(&feieyun.Text{Content: fmt.Sprintf("桌號: %s", bill.TableLabel)})
		for _, order := range orderNumbers {
			pc.AddLines(
				&feieyun.Text{
					Content: fmt.Sprintf("%s %.2fX%d", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)})
			for _, option := range order.Order.Specification {
				pc.AddLines(
					&feieyun.Text{
						Content: fmt.Sprintf("- %s +%.2f", option.R, float64(order.Order.Extra(option))/100)})
			}
		}
		pc.AddLines(
			&feieyun.Text{
				Content: fmt.Sprintf("合計: %2.f元", float64(bill.Total()/100))})
		pc.AddLines(
			&feieyun.Text{
				Content: fmt.Sprintf("時間: %s", bill.CreatedAt.Add(time.Hour*8).Format("2006-01-02 15:04"))})
	}
	pc.AddDiv(width)
	pc.AddLines(
		&feieyun.Bold{Content: &feieyun.Text{
			Content: fmt.Sprintf("總合計: %2.f元",
				math.Floor((float64(total) / 100 * _offset)))}})
	return pc.String()
}
