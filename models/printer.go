package models

import (
	"fmt"
	"math"
	"time"

	abstract "github.com/Dparty/dao/abstract"
	restaurantDao "github.com/Dparty/dao/restaurant"
)

func NewPrinter(entity restaurantDao.Printer) *Printer {
	return &Printer{entity: entity}
}

type Printer struct {
	entity restaurantDao.Printer
}

func (p Printer) ID() uint {
	return p.entity.ID()
}

func (p Printer) Sn() string {
	return p.entity.Sn
}

func (p Printer) Name() string {
	return p.entity.Name
}

func (p Printer) Description() string {
	return p.entity.Description
}

func (p Printer) Type() string {
	return string(p.entity.Type)
}

func (p Printer) Delete() bool {
	items := itemRepository.List("restaurant_id = ?", p.Owner().ID())
	for _, item := range items {
		for _, printer := range item.Printers {
			if printer == p.ID() {
				return false
			}
		}
	}
	printerRepository.Delete(p.ID())
	return true
}

func (p Printer) SetOwner(r abstract.Owner) *Printer {
	return &p
}

func (p Printer) Owner() *restaurantDao.Restaurant {
	return restaurantRepository.GetById(p.entity.Owner().ID())
}

type OrderNumber struct {
	Number int
	Order  restaurantDao.Order
}

func PrintBill(printers []restaurantDao.Printer, restaurantName string, bill restaurantDao.Bill, table restaurantDao.Table, offset int64, reprint bool) {
	timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	orderNumbers := make([]OrderNumber, 0)
	for _, order := range bill.Orders {
		orderNumbers = PrintHelper(order, orderNumbers)
	}
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurantName)
	content += fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
	content += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
	content += "--------------------------------<BR>"
	for _, order := range orderNumbers {
		content += fmt.Sprintf("<B>%s %.2fX%d</B><BR>", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)
		attributes := ""
		for _, option := range order.Order.Specification {
			attributes += fmt.Sprintf("<B>|-- %s +%.2f</B><BR>", option.Right, float64(order.Order.Extra(option))/100)
		}
		content += attributes
	}
	for _, order := range orderNumbers {
		a := fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
		a += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
		a += fmt.Sprintf("<B>%s X%d</B><BR>", order.Order.Item.Name, order.Number)
		for _, option := range order.Order.Specification {
			a += fmt.Sprintf("<B>|--  %s</B><BR>", option.Right)
		}
		for _, printer := range order.Order.Item.Printers {
			foodPrinter := GetPrinter(printer)
			p, _ := printerFactory.Connect(foodPrinter.Sn)
			p.Print(a+"<BR>"+timestring, "")
		}
	}
	content += "--------------------------------<BR>"
	_offset := float64(offset+100) / 100
	content += fmt.Sprintf("<B>合計: %.2f元</B><BR>", math.Floor(float64(bill.Total())/100*_offset))
	content += timestring
	for _, printer := range printers {
		if printer.Type == "BILL" {
			p, _ := printerFactory.Connect(printer.Sn)
			p.Print(content, "")
		}
	}
}

func GetPrinter(id uint) restaurantDao.Printer {
	var printer restaurantDao.Printer
	db.Where("id = ?", id).Find(&printer)
	return printer
}

func PrintHelper(order restaurantDao.Order, orders []OrderNumber) []OrderNumber {
	for i, o := range orders {
		if order.Equal(o.Order) {
			orders[i].Number++
			return orders
		}
	}
	orders = append(orders, OrderNumber{
		Number: 1,
		Order:  order,
	})
	return orders
}
