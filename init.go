package restaurantservice

import (
	"github.com/Dparty/common/cloud"
	"github.com/Dparty/common/config"
	"github.com/Dparty/dao"
	"github.com/Dparty/dao/restaurant"
	"github.com/Dparty/feieyun"
)

var db = dao.GetDBInstance()
var CosClient cloud.CosClient
var Bucket string
var printerFactory = feieyun.NewPrinterFactory(config.GetString("feieyun.user"), config.GetString("feieyun.ukey"), config.GetString("feieyun.url"))
var restaurantRepository = restaurant.GetRestaurantRepository()
var itemRepository = restaurant.GetItemRepository()
var printerRepository = restaurant.GetPrinterRepository()
var tableRepository = restaurant.GetTableRepository()
var billRepository = restaurant.GetBillRepository()
