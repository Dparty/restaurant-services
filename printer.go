package restaurantservice

import (
	"fmt"
	"math"
	"time"

	abstract "github.com/Dparty/common/abstract"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
)

func NewPrinter(entity restaurantDao.Printer) *Printer {
	return &Printer{entity: entity}
}

type Printer struct {
	entity restaurantDao.Printer
}

func (p Printer) Print(printContent feieyun.PrintContent, reprint bool) {
	printer, _ := printerFactory.Connect(p.Sn())
	printer.Print(printContent.String(), "")
}

func (p Printer) PrintBill(restaurantName string, bill restaurantDao.Bill, table restaurantDao.Table, offset int64, reprint bool) {
	timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	orderNumbers := make([]OrderNumber, 0)
	for _, order := range bill.Orders {
		orderNumbers = MargeOrderHelper(order, orderNumbers)
	}
	var printContent feieyun.PrintContent
	printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: restaurantName}})
	printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("餐號: %d", bill.PickUpCode)}})
	printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("桌號: %s", table.Label)}})
	printContent.AddDiv(p.Width())
	for _, order := range orderNumbers {
		printContent.AddLines(
			&feieyun.Bold{
				Content: &feieyun.Text{
					Content: fmt.Sprintf("%s %.2fX%d", order.Order.Item.Name, float64(order.Order.Item.Pricing)/100, order.Number)}})
		for _, option := range order.Order.Specification {
			printContent.AddLines(
				&feieyun.Bold{
					Content: &feieyun.Text{
						Content: fmt.Sprintf("- %s +%.2fX%d", option.R, float64(order.Order.Extra(option))/100, order.Number)}})
		}
	}
	printContent.AddDiv(p.Width())
	printContent.AddLines(&feieyun.Bold{Content: &feieyun.Text{Content: fmt.Sprintf("合計: %.2f元", math.Floor(float64(bill.Total())/100))}})
	printContent.AddLines(&feieyun.Text{Content: timestring})
	p.Print(printContent, reprint)
}

func (p Printer) Width() int {
	if p.Model() == "58mm" {
		return 32
	}
	return 48
}

func (p Printer) ID() uint {
	return p.entity.ID()
}

func (p Printer) Sn() string {
	return p.entity.Sn
}

func (p *Printer) SetSn(sn string) *Printer {
	p.entity.Sn = sn
	return p
}

func (p Printer) Name() string {
	return p.entity.Name
}

func (p *Printer) SetName(name string) *Printer {
	p.entity.Name = name
	return p
}

func (p Printer) Description() string {
	return p.entity.Description
}

func (p *Printer) SetDescription(description string) *Printer {
	p.entity.Description = description
	return p
}

func (p Printer) Type() string {
	return string(p.entity.Type)
}

func (p *Printer) SetType(t string) *Printer {
	p.entity.Type = restaurantDao.PrinterType(t)
	return p
}

func (p *Printer) SetModel(model string) *Printer {
	p.entity.PrinterModel = model
	return p
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

func (p *Printer) Submit() {
	printerRepository.Save(&p.entity)
}

func (p *Printer) SetOwner(r abstract.Owner) *Printer {
	p.entity.SetOwner(r)
	return p
}

func (p Printer) Owner() *restaurantDao.Restaurant {
	return restaurantRepository.GetById(p.entity.Owner().ID())
}

type OrderNumber struct {
	Number int
	Order  restaurantDao.Order
}

func (p Printer) Model() string {
	if p.entity.PrinterModel == "" {
		return "58mm"
	}
	return p.entity.PrinterModel
}

func PrintBill(printers []Printer, restaurantName string, bill restaurantDao.Bill, table restaurantDao.Table, offset int64, reprint bool) {
	timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	orderNumbers := make([]OrderNumber, 0)
	for _, order := range bill.Orders {
		orderNumbers = MargeOrderHelper(order, orderNumbers)
	}
	for _, order := range orderNumbers {
		var pc feieyun.PrintContent
		pc.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("餐號: %d", bill.PickUpCode)}},
			&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("桌號: %s", table.Label)}},
			&feieyun.Bold{Content: &feieyun.Text{Content: fmt.Sprintf("%s X%d", order.Order.Item.Name, order.Number)}})
		for _, option := range order.Order.Specification {
			pc.AddLines(&feieyun.Bold{Content: &feieyun.Text{Content: fmt.Sprintf("-  %sX%d", option.R, order.Number)}})
		}
		pc.AddLines(&feieyun.Text{Content: timestring})
		for _, printer := range order.Order.Item.Printers {
			p, err := printerService.GetById(printer)
			if err == nil {
				p.Print(pc, reprint)
			}
		}
	}
	for _, printer := range printers {
		if printer.Type() == "BILL" {
			printer.PrintBill(restaurantName, bill, table, offset, reprint)
		}
	}
}

func MargeOrderHelper(order restaurantDao.Order, orders []OrderNumber) []OrderNumber {
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
