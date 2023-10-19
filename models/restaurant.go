package models

import (
	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
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

// SetOwner implements interfaces.Asset.
func (r *Restaurant) SetOwner(owner abstract.Owner) *Restaurant {
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
	// TODO: Check if printers doesn't exist.
	item := restaurant.Item{
		RestaurantId: r.ID(),
		Pricing:      pricing,
		Attributes:   attributes,
		Images:       images,
		Tags:         tags,
		Printers:     printers,
	}
	itemRepository.Save(&item)
	return NewItem(item), nil
}

func (r Restaurant) CreateTable(label string, x, y int64) (Table, error) {
	// TODO: Check if table label exist or position conflict
	table := restaurant.Table{
		RestaurantId: r.ID(),
		Label:        label,
		X:            x,
		Y:            y,
	}
	tableRepository.Save(&table)
	return Table{table}, nil
}

type Table struct {
	entity restaurant.Table
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
