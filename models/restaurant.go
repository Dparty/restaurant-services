package models

import (
	"github.com/Dparty/common/fault"
	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
	"github.com/chenyunda218/golambda"
)

func NewRestaurant(entity restaurant.Restaurant) Restaurant {
	return Restaurant{entity: entity}
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
	for _, table := range r.entity.Tables() {
		tables = append(tables, Table{table})
	}
	return tables
}

func (r Restaurant) Items() []Item {
	var items []Item
	for _, item := range r.entity.Items() {
		items = append(items, NewItem(item))
	}
	return items
}

func (r Restaurant) Printers() []Printer {
	var printers []Printer
	for _, printer := range r.entity.Printers() {
		printers = append(printers, Printer{printer})
	}
	return printers
}

func (r *Restaurant) Update(name, description string) *Restaurant {
	return r
}

func (r Restaurant) CreateItem(name string, pricing int64, attributes restaurant.Attributes, images, tags []string, printers []uint) (Item, error) {
	item, err := itemRepository.Save(&restaurant.Item{
		RestaurantId: r.ID(),
		Pricing:      pricing,
		Attributes:   attributes,
		Images:       images,
		Tags:         tags,
		Printers:     printers,
	})
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

func (r Restaurant) CreatePrinter(t, sn, name, description string) (Printer, error) {
	printer := restaurant.Printer{
		RestaurantId: r.ID(),
		Name:         name,
		Sn:           sn,
		Type:         restaurant.PrinterType(t),
		Description:  description,
	}
	printerRepository.Save(&printer)
	return Printer{printer}, nil
}

func NewTable(entity restaurant.Table) *Table {
	return &Table{entity: entity}
}

type Table struct {
	entity restaurant.Table
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

func (t Table) Entity() restaurant.Table {
	return t.entity
}

func (t Table) Delete() bool {
	return tableRepository.Delete(&t.entity).RowsAffected != 0
}
