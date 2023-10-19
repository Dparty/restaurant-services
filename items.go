package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/restaurant-services/models"
	"gorm.io/gorm"
)

func NewItemService(inject *gorm.DB) ItemService {
	return ItemService{itemRepository: restaurantDao.NewItemRepository(inject)}
}

type ItemService struct {
	itemRepository restaurantDao.ItemRepository
}

func (i ItemService) GetById(id uint) (models.Item, error) {
	entity := i.itemRepository.GetById(id)
	if entity == nil {
		return models.Item{}, fault.ErrNotFound
	}
	return models.NewItem(*entity), nil
}
