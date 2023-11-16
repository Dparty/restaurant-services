package models

import (
	"fmt"
	"math"

	abstract "github.com/Dparty/common/abstract"
	restaurantDao "github.com/Dparty/dao/restaurant"
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
	restaurant := restaurantRepository.GetById(t.Owner().ID())
	printers := restaurant.Printers()
	status := "SUBMITTED"
	bills := t.Bills(&status)
	if len(bills) == 0 {
		return
	}
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurant.Name)
	content += fmt.Sprintf("<CB>桌號: %s</CB><BR>", t.Label())
	content += FinishString(
		offset,
		golambda.Map(bills,
			func(_ int, b Bill) restaurantDao.Bill {
				return b.Entity()
			}))
	for _, printer := range printers {
		if printer.Type == "BILL" {
			p, _ := printerFactory.Connect(printer.Sn)
			p.Print(content, "")
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

func FinishString(offset int64, bills []restaurantDao.Bill) string {
	_offset := float64(offset+100) / 100
	content := ""
	total := 0
	for _, bill := range bills {
		total += int(bill.Total())
		orderNumbers := make([]OrderNumber, 0)
		for _, order := range bill.Orders {
			orderNumbers = PrintHelper(order, orderNumbers)
		}
		content += "--------------------------------<BR>"
		content += fmt.Sprintf("餐號: %d<BR>", bill.PickUpCode)
		content += fmt.Sprintf("桌號: %s<BR>", bill.TableLabel)
		for _, order := range orderNumbers {
			content += fmt.Sprintf("%s %.2fX%d<BR>", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)
			for _, option := range order.Order.Specification {
				content += fmt.Sprintf("|- %s +%.2f<BR>", option.Right, float64(order.Order.Extra(option))/100)
			}
		}
		content += fmt.Sprintf("合計: %2.f元<BR>", float64(bill.Total()/100))
	}
	content += "--------------------------------<BR>"
	content += fmt.Sprintf("<B>總合計: %2.f元</B><BR>",
		math.Floor((float64(total) / 100 * _offset)))
	return content
}
