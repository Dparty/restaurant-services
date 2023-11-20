package restaurantservice

import (
	"github.com/Dparty/common/cloud"
	"github.com/Dparty/common/config"
	"github.com/Dparty/dao/auth"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/Dparty/restaurant-services/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

var CosClient cloud.CosClient
var Bucket string
var printerFactory feieyun.PrinterFactory
var accountRepository auth.AccountRepository
var restaurantRepository restaurant.RestaurantRepository
var itemRepository restaurant.ItemRepository
var printerRepository restaurant.PrinterRepository

func Init(inject *gorm.DB) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.host") + ":6379",
		Password: "",
		DB:       0,
	})
	db = inject
	auth.Init(inject)
	restaurant.Init(inject)
	models.Init(inject)
	accountRepository = auth.NewAccountRepository(inject)
	restaurantRepository = restaurant.NewRestaurantRepository(inject)
	itemRepository = restaurant.NewItemRepository(inject)
	printerRepository = restaurant.NewPrinterRepository(inject)
	Bucket = config.GetString("cos.Bucket")
	CosClient.Region = config.GetString("cos.Region")
	CosClient.SecretID = config.GetString("cos.SecretID")
	CosClient.SecretKey = config.GetString("cos.SecretKey")
	printerFactory = feieyun.NewPrinterFactory(config.GetString("feieyun.user"), config.GetString("feieyun.ukey"), config.GetString("feieyun.url"))
}
