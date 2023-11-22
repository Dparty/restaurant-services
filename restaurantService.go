package restaurantservice

import (
	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"gorm.io/gorm"
)

func NewRestaurantService(inject *gorm.DB) *RestaurantService {
	return &RestaurantService{restaurantRepository: restaurantDao.NewRestaurantRepository(inject)}
}

type RestaurantService struct {
	restaurantRepository restaurantDao.RestaurantRepository
}

func (r RestaurantService) CreateRestaurant(owner abstract.Owner, name, description string) Restaurant {
	entity := r.restaurantRepository.Create(owner, name, description)
	return NewRestaurant(entity)
}

func (r RestaurantService) UpdateRestaurant(id uint, name, description string, categories []string) (Restaurant, error) {
	restaurant, err := r.GetRestaurant(id)
	if err != nil {
		return restaurant, err
	}
	restaurant.SetName(name).SetDescription(description).SetCategories(categories).Submit()
	return restaurant, err
}

func (r RestaurantService) GetRestaurant(id uint) (Restaurant, error) {
	entity := r.restaurantRepository.GetById(id)
	if entity == nil {
		return Restaurant{}, fault.ErrNotFound
	}
	return NewRestaurant(*entity), nil
}

func (r RestaurantService) List(ownerId *uint) []Restaurant {
	var restaurants []Restaurant
	for _, r := range r.restaurantRepository.ListBy(ownerId) {
		restaurants = append(restaurants, NewRestaurant(r))
	}
	return restaurants
}
