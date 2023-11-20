package restaurantservice

import (
	"github.com/Dparty/common/fault"
	restaurantDao "github.com/Dparty/dao/restaurant"

	"gorm.io/gorm"
)

func NewPrinterService(inject *gorm.DB) PrinterService {
	return PrinterService{printerRepository: restaurantDao.NewPrinterRepository(inject)}
}

type PrinterService struct {
	printerRepository restaurantDao.PrinterRepository
}

func (p PrinterService) GetById(id uint) (*Printer, error) {
	entity := p.printerRepository.GetById(id)
	if entity == nil {
		return nil, fault.ErrNotFound
	}
	return NewPrinter(*entity), nil
}
