package restaurantservice

import (
	"github.com/Dparty/common/cloud"
	"github.com/Dparty/common/config"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"gorm.io/gorm"
)

var db *gorm.DB
var CosClient cloud.CosClient
var Bucket string
var printerFactory feieyun.PrinterFactory
var restaurantRepository *restaurant.RestaurantRepository
var itemRepository *restaurant.ItemRepository
var printerRepository *restaurant.PrinterRepository
var tableRepository *restaurant.TableRepository
var billRepository *restaurant.BillRepository

func Init(inject *gorm.DB) {
	db = inject
	billRepository = restaurant.GetBillRepository()
	tableRepository = restaurant.GetTableRepository()
	restaurantRepository = restaurant.GetRestaurantRepository()
	itemRepository = restaurant.GetItemRepository()
	printerRepository = restaurant.GetPrinterRepository()
	Bucket = config.GetString("cos.Bucket")
	CosClient.Region = config.GetString("cos.Region")
	CosClient.SecretID = config.GetString("cos.SecretID")
	CosClient.SecretKey = config.GetString("cos.SecretKey")
	printerFactory = feieyun.NewPrinterFactory(config.GetString("feieyun.user"), config.GetString("feieyun.ukey"), config.GetString("feieyun.url"))
	restaurantService = NewRestaurantService()
}
