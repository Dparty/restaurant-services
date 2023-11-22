package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
)

var itemService *ItemService

func GetItemService() *ItemService {
	if itemService == nil {
		itemService = NewItemService()
	}
	return itemService
}

func NewItemService() *ItemService {
	return &ItemService{itemRepository: restaurantDao.GetItemRepository()}
}

type ItemService struct {
	itemRepository *restaurantDao.ItemRepository
}

func (i ItemService) GetById(id uint) (Item, error) {
	entity := i.itemRepository.GetById(id)
	if entity == nil {
		return Item{}, fault.ErrNotFound
	}
	return NewItem(*entity), nil
}
