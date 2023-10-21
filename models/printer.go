package models

import (
	"fmt"

	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
)

func NewPrinter(entity restaurant.Printer) *Printer {
	return &Printer{entity: entity}
}

type Printer struct {
	entity restaurant.Printer
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

func (p Printer) Owner() *restaurant.Restaurant {
	return restaurantRepository.GetById(p.entity.Owner().ID())
}

type OrderNumber struct {
	Number int
	Order  restaurant.Order
}

func PrintBill(printers []restaurant.Printer, restaurantName string, bill restaurant.Bill, table restaurant.Table, reprint bool) {
	orderNumbers := make([]OrderNumber, 0)
	for _, order := range bill.Orders {
		orderNumbers = PrintHelper(order, orderNumbers)
	}
	content := ""
	content += fmt.Sprintf("<CB>%s</CB><BR>", restaurantName)
	content += fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
	content += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
	content += "--------------------------------<BR>"
	var printersString map[uint]string = make(map[uint]string)
	for _, order := range orderNumbers {
		content += fmt.Sprintf("<B>%s %.2f X %d</B><BR>", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)
		attributes := ""
		attributesWithoutMonth := ""
		for _, option := range order.Order.Specification {
			attributes += fmt.Sprintf("<B>|-- %s +%.2f</B><BR>", option.Right, float64(order.Order.Extra(option))/100)
			attributesWithoutMonth += fmt.Sprintf("<B>|--  %s</B><BR>", option.Right)
		}
		content += attributes
		for _, printer := range order.Order.Item.Printers {
			_, ok := printersString[printer]
			if !ok {
				printersString[printer] = fmt.Sprintf("<CB>餐號: %d</CB><BR>", bill.PickUpCode)
				printersString[printer] += fmt.Sprintf("<CB>桌號: %s</CB><BR>", table.Label)
			}
			printersString[printer] += fmt.Sprintf("<B>%s X %d</B><BR>", order.Order.Item.Name, order.Number)
			printersString[printer] += attributesWithoutMonth
		}
	}
	for k, v := range printersString {
		fmt.Println(k, v)
		// foodPrinter := printerRepository.GetById(k)
		// p, _ := printerFactory.Connect(foodPrinter.Sn)
		// p.Print(v, "")
	}
	content += "--------------------------------<BR>"
	content += fmt.Sprintf("合計: %.2f元<BR>", float64(bill.Total())/100)
	// for _, printer := range printers {
	// 	if printer.Type == "BILL" {
	// 		p, _ := printerFactory.Connect(printer.Sn)
	// 		p.Print(content, "")
	// 	}
	// }
}
func PrintHelper(order restaurant.Order, orders []OrderNumber) []OrderNumber {
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
