package models

import (
	abstract "github.com/Dparty/dao/abstract"
	"github.com/Dparty/dao/restaurant"
)

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

func (t Table) Bills(status *string) []Bill {
	var bills []Bill
	for _, b := range t.entity.Bills(status) {
		bills = append(bills, NewBill(b))
	}
	return bills
}

func (t Table) Delete() bool {
	return tableRepository.Delete(&t.entity).RowsAffected != 0
}
