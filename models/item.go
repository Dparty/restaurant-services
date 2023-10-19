package models

import (
	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
)

func NewItem(entity restaurant.Item) Item {
	return Item{entity: entity}
}

type Item struct {
	entity restaurant.Item
}

func (i Item) ID() uint {
	return i.entity.ID()
}

func (i Item) Entity() restaurant.Item {
	return i.entity
}

func (i *Item) Update(name string, pricing int64, attributes restaurant.Attributes, images, tags []string, printers []uint) (*Item, error) {
	return i, nil
}

func (i Item) Delete() bool {
	ctx := itemRepository.Delete(&i.entity)
	return ctx.RowsAffected != 0
}

func (i Item) Owner() abstract.Owner {
	return i.entity.Owner()
}
