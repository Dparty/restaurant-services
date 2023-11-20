package restaurantservice

import (
	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
	"github.com/Dparty/restaurant-services/models"
	"gorm.io/gorm"
)

func NewRestaurantService(inject *gorm.DB) RestaurantService {
	return RestaurantService{restaurantRepository: restaurantDao.NewRestaurantRepository(inject)}
}

type RestaurantService struct {
	restaurantRepository restaurantDao.RestaurantRepository
}

func (r RestaurantService) CreateRestaurant(owner abstract.Owner, name, description string) models.Restaurant {
	entity := r.restaurantRepository.Create(owner, name, description)
	return models.NewRestaurant(entity)
}

func (r RestaurantService) UpdateRestaurant(id uint, name, description string, categories []string) (models.Restaurant, error) {
	restaurant, err := r.GetRestaurant(id)
	if err != nil {
		return restaurant, err
	}
	entity := restaurantRepository.GetById(id)
	entity.Name = name
	entity.Description = description
	entity.Categories = categories

	restaurantRepository.Save(entity)
	n := models.NewRestaurant(*entity)
	return n, err
}

func (r RestaurantService) GetRestaurant(id uint) (models.Restaurant, error) {
	entity := r.restaurantRepository.GetById(id)
	if entity == nil {
		return models.Restaurant{}, fault.ErrNotFound
	}
	return models.NewRestaurant(*entity), nil
}

func (r RestaurantService) List(ownerId *uint) []models.Restaurant {
	var restaurants []models.Restaurant
	for _, r := range r.restaurantRepository.ListBy(ownerId) {
		restaurants = append(restaurants, models.NewRestaurant(r))
	}
	return restaurants
}
