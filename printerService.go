package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"
)

var printerService *PrinterService

func GetPrinterService() *PrinterService {
	if printerService == nil {
		printerService = NewPrinterService()
	}
	return printerService
}

func NewPrinterService() *PrinterService {
	return &PrinterService{printerRepository: restaurantDao.GetPrinterRepository()}
}

type PrinterService struct {
	printerRepository *restaurantDao.PrinterRepository
}

func (p PrinterService) GetById(id uint) (*Printer, error) {
	entity := p.printerRepository.GetById(id)
	if entity == nil {
		return nil, fault.ErrNotFound
	}
	return NewPrinter(*entity), nil
}
