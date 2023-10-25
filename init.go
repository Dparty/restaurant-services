package restaurantservice

import (
	"fmt"

	"github.com/Dparty/common/cloud"
	"github.com/Dparty/dao/auth"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/Dparty/restaurant-services/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var CosClient cloud.CosClient
var Bucket string
var printerFactory feieyun.PrinterFactory
var accountRepository auth.AccountRepository
var restaurantRepository restaurant.RestaurantRepository
var itemRepository restaurant.ItemRepository
var printerRepository restaurant.PrinterRepository

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	Bucket = viper.GetString("cos.Bucket")
	CosClient.Region = viper.GetString("cos.Region")
	CosClient.SecretID = viper.GetString("cos.SecretID")
	CosClient.SecretKey = viper.GetString("cos.SecretKey")
}

func Init(inject *gorm.DB) {
	db = inject
	auth.Init(inject)
	restaurant.Init(inject)
	models.Init(inject)
	accountRepository = auth.NewAccountRepository(inject)
	restaurantRepository = restaurant.NewRestaurantRepository(inject)
	itemRepository = restaurant.NewItemRepository(inject)
	printerRepository = restaurant.NewPrinterRepository(inject)
}
