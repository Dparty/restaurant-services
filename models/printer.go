package models

import (
	"fmt"

	"github.com/Dparty/dao/restaurant"
)

type Printer struct {
	entity restaurant.Printer
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
		foodPrinter := printerRepository.GetById(k)
		p, _ := printerFactory.Connect(foodPrinter.Sn)
		p.Print(v, "")
	}
	content += "--------------------------------<BR>"
	content += fmt.Sprintf("合計: %.2f元<BR>", float64(bill.Total())/100)
	for _, printer := range printers {
		if printer.Type == "BILL" {
			p, _ := printerFactory.Connect(printer.Sn)
			p.Print(content, "")
		}
	}
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
