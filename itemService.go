package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"gorm.io/gorm"
)

func NewItemService(inject *gorm.DB) ItemService {
	return ItemService{itemRepository: restaurantDao.NewItemRepository(inject)}
}

type ItemService struct {
	itemRepository restaurantDao.ItemRepository
}

func (i ItemService) GetById(id uint) (Item, error) {
	entity := i.itemRepository.GetById(id)
	if entity == nil {
		return Item{}, fault.ErrNotFound
	}
	return NewItem(*entity), nil
}
