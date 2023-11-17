package restaurantservice

import (
	"github.com/Dparty/common/cloud"
	"github.com/Dparty/dao/auth"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/Dparty/restaurant-services/models"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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

func Init(v *viper.Viper, inject *gorm.DB) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     v.GetString("redis.host") + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	db = inject
	auth.Init(inject)
	restaurant.Init(inject)
	models.Init(inject)
	accountRepository = auth.NewAccountRepository(inject)
	restaurantRepository = restaurant.NewRestaurantRepository(inject)
	itemRepository = restaurant.NewItemRepository(inject)
	printerRepository = restaurant.NewPrinterRepository(inject)
	Bucket = v.GetString("cos.Bucket")
	CosClient.Region = v.GetString("cos.Region")
	CosClient.SecretID = v.GetString("cos.SecretID")
	CosClient.SecretKey = v.GetString("cos.SecretKey")
	printerFactory = feieyun.NewPrinterFactory(v.GetString("feieyun.user"), v.GetString("feieyun.ukey"), v.GetString("feieyun.url"))
}
