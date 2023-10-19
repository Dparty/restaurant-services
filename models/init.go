package models

import (
	"fmt"

	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var itemRepository restaurant.ItemRepository
var tableRepository restaurant.TableRepository
var printerRepository restaurant.PrinterRepository
var printerFactory feieyun.PrinterFactory

func Init(inject *gorm.DB) {
	restaurant.Init(inject)
	itemRepository = restaurant.NewItemRepository(inject)
	tableRepository = restaurant.NewTableRepository(inject)
	printerRepository = restaurant.NewPrinterRepository(inject)
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
