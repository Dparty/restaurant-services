package models

import (
	"fmt"

	"github.com/Dparty/common/cloud"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var restaurantRepository restaurant.RestaurantRepository
var itemRepository restaurant.ItemRepository
var tableRepository restaurant.TableRepository
var printerRepository restaurant.PrinterRepository
var billRepository restaurant.BillRepository
var printerFactory feieyun.PrinterFactory

var CosClient cloud.CosClient
var Bucket string

func Init(inject *gorm.DB) {
	db = inject
	restaurant.Init(inject)
	itemRepository = restaurant.NewItemRepository(inject)
	tableRepository = restaurant.NewTableRepository(inject)
	printerRepository = restaurant.NewPrinterRepository(inject)
	restaurantRepository = restaurant.NewRestaurantRepository(inject)
	billRepository = restaurant.NewBillRepository(inject)
}

func init() {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	user := viper.GetString("feieyun.user")
	ukey := viper.GetString("feieyun.ukey")
	url := viper.GetString("feieyun.url")
	printerFactory = feieyun.NewPrinterFactory(user, ukey, url)
}

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
