package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"gorm.io/gorm"
)

var tableService *TableService

func GetTableService() *TableService {
	if tableService == nil {
		tableService = NewTableService(nil)
	}
	return tableService
}

func NewTableService(inject *gorm.DB) *TableService {
	return &TableService{tableRepository: restaurantDao.GetTableRepository()}
}

type TableService struct {
	tableRepository *restaurantDao.TableRepository
}

func (t TableService) GetById(id uint) (*Table, error) {
	table := t.tableRepository.Find(id)
	if table == nil {
		return nil, fault.ErrNotFound
	}
	return NewTable(*table), nil
}
