package models

import (
	"fmt"

	abstract "github.com/Dparty/dao/abstract"
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

func (t Table) Finish() {
	restaurant := restaurantRepository.GetById(t.Owner().ID())
	printers := restaurant.Printers()
	fmt.Println(printers)
	status := "SUBMIT"
	bills := t.Bills(&status)
	for _, bill := range bills {
		bill.Finish()
	}
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurant.Name)
	content += fmt.Sprintf("<CB>桌號: %s</CB><BR>", t.Label())
	content += FinishString(golambda.Map(bills, func(_ int, b Bill) restaurantDao.Bill {
		return b.Entity()
	}))
	for _, printer := range printers {
		if printer.Type == "BILL" {
			p, _ := printerFactory.Connect(printer.Sn)
			p.Print(content, "")
		}
	}
}

func FinishString(bills []restaurantDao.Bill) string {
	content := ""
	for _, bill := range bills {
		orderNumbers := make([]OrderNumber, 0)
		for _, order := range bill.Orders {
			orderNumbers = PrintHelper(order, orderNumbers)
		}
		content += "--------------------------------<BR>"
		for _, order := range orderNumbers {
			content += fmt.Sprintf("%sX%d<BR>", order.Order.Item.Name, order.Number)
			for _, option := range order.Order.Specification {
				content += fmt.Sprintf("|- %s<BR>", option.Right)
			}
		}
		content += fmt.Sprintf("合計: %2.f元<BR>", float64(bill.Total()/100))
	}
	return content
}
