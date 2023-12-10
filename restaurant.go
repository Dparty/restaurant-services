package restaurantservice

import (
	"time"

	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/dao/restaurant"
	"github.com/chenyunda218/golambda"
)

func NewRestaurant(entity restaurant.Restaurant) *Restaurant {
	return &Restaurant{entity: entity}
}

type Restaurant struct {
	entity restaurant.Restaurant
}

// Owner implements interfaces.Asset.
func (r *Restaurant) Owner() abstract.Owner {
	return r.entity.Owner()
}

func (r Restaurant) Own(asset abstract.Asset) bool {
	return r.ID() == asset.Owner().ID()
}

// SetOwner implements interfaces.Asset.
func (r *Restaurant) SetOwner(owner abstract.Owner) abstract.Asset {
	r.entity.SetOwner(owner)
	return r
}

func (r *Restaurant) Delete() {
	for _, item := range r.Items() {
		item.Delete()
	}
	for _, table := range r.Tables() {
		table.Delete()
	}
	for _, printer := range r.Printers() {
		printer.Delete()
	}
	restaurantRepository.Delete(&r.entity)
}

func (r Restaurant) ID() uint {
	return r.entity.ID()
}

func (r Restaurant) Name() string {
	return r.entity.Name
}

func (r Restaurant) Description() string {
	return r.entity.Description
}

func (r Restaurant) Entity() restaurant.Restaurant {
	return r.entity
}

func (r Restaurant) Tables() []Table {
	var tables []Table
	for _, table := range tableRepository.List("restaurant_id = ?", r.ID()) {
		tables = append(tables, Table{table})
	}
	return tables
}

func (r Restaurant) Items() []Item {
	var items []Item
	for _, item := range itemRepository.List("restaurant_id = ?", r.ID()) {
		items = append(items, NewItem(item))
	}
	return items
}

func (r Restaurant) Printers() []Printer {
	var printers []Printer
	for _, printer := range printerRepository.List("restaurant_id = ?", r.ID()) {
		printers = append(printers, Printer{printer})
	}
	return printers
}

func (r Restaurant) PickUpCode() int64 {
	return r.entity.PickUpCode()
}

func (r *Restaurant) SetName(name string) *Restaurant {
	r.entity.Name = name
	return r
}

func (r *Restaurant) SetDescription(description string) *Restaurant {
	r.entity.Description = description
	return r
}

func (r *Restaurant) AddPrinter(printer *Printer) {

}

func (r *Restaurant) SetCategories(categories []string) *Restaurant {
	r.entity.Categories = categories
	return r
}

func removeDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (r *Restaurant) Categories() []string {
	var categories []string
	categories = append(categories, r.entity.Categories...)
	for _, item := range r.Items() {
		categories = append(categories, item.Categories()...)
	}
	return removeDuplicateElement(categories)
}

func (r *Restaurant) Submit() *Restaurant {
	db.Where("id = ?", r.ID()).Updates(&r.entity)
	return r
}

func (r Restaurant) CreateItem(name string, pricing int64, attributes restaurant.Attributes, images, tags []string, printers []uint, status string) (Item, error) {
	i := restaurant.Item{
		Name:         name,
		RestaurantId: r.ID(),
		Pricing:      pricing,
		Attributes:   attributes,
		Images:       images,
		Tags:         tags,
		Printers:     printers,
		Status:       status,
	}
	item, err := itemRepository.Save(&i)
	if err != nil {
		return Item{}, err
	}
	return NewItem(*item), nil
}

func (r Restaurant) CreateTable(label string, x, y int64) (Table, error) {
	if len(golambda.Filter(r.Tables(), func(_ int, table Table) bool {
		return table.Label() == label || (x == table.X() && y == table.Y())
	})) != 0 {
		return Table{}, fault.ErrCreateTableConflict
	}
	table := restaurant.Table{
		RestaurantId: r.ID(),
		Label:        label,
		X:            x,
		Y:            y,
	}
	tableRepository.Save(&table)
	return Table{table}, nil
}

func (r Restaurant) CreatePrinter(t, sn, name, description, model string) (Printer, error) {
	printer := restaurant.Printer{
		RestaurantId: r.ID(),
		Name:         name,
		Sn:           sn,
		Type:         restaurant.PrinterType(t),
		Description:  description,
		PrinterModel: model,
	}
	printerRepository.Save(&printer)
	return Printer{printer}, nil
}

func NewTable(entity restaurant.Table) *Table {
	return &Table{entity: entity}
}

func (r Restaurant) ListBills(restaurantId uint,
	tableId *uint,
	status *string,
	startAt, endAt *time.Time) []Bill {
	ctx := db.Model(&restaurant.Bill{})
	ctx = ctx.Where("restaurant_id = ?", restaurantId)
	if tableId != nil {
		ctx = ctx.Where("table_id = ?", *tableId)
	}
	if status != nil {
		ctx = ctx.Where("status = ?", *tableId)
	}
	if startAt != nil {
		ctx = ctx.Where("created_at >= ?", *startAt)
	}
	if endAt != nil {
		ctx = ctx.Where("created_at <= ?", *endAt)
	}
	var bs []restaurant.Bill
	ctx.Find(&bs)
	var bills []Bill
	for _, b := range bs {
		bills = append(bills, NewBill(b))
	}
	return bills
}

func (r Restaurant) CreateDiscount(label string, offset int64) Discount {
	discount := restaurant.Discount{
		RestaurantId: r.ID(),
		Label:        label,
		Offset:       offset,
	}
	discountRepository.Save(&discount)
	return NewDiscount(discount)
}

func (r Restaurant) Discounts() []Discount {
	var discounts []Discount
	for _, discount := range discountRepository.ListBy(r.ID()) {
		discounts = append(discounts, NewDiscount(discount))
	}
	return discounts
}
