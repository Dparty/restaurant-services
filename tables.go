package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/restaurant-services/models"
	"gorm.io/gorm"
)

func NewTableService(inject *gorm.DB) TableService {
	return TableService{tableRepository: restaurantDao.NewTableRepository(inject)}
}

type TableService struct {
	tableRepository restaurantDao.TableRepository
}

func (t TableService) GetById(id uint) (*models.Table, error) {
	table := t.tableRepository.Find(id)
	if table == nil {
		return nil, fault.ErrNotFound
	}
	return models.NewTable(*table), nil
}
