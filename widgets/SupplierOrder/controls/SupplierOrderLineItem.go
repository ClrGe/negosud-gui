package controls

import (
	"fyne.io/fyne/v2/widget"
)

type SupplierOrderLineItem struct {
	BottleId          int
	StorageLocationId int
	SelectBottle      *widget.Select
	EntryQuantity     *widget.Entry
	BottleMap         map[string]int
}

func (i *SupplierOrderLineItem) Bind(bottleId int, storageLocationId int) {
	i.BottleId = bottleId
	i.StorageLocationId = storageLocationId
}

func (i *SupplierOrderLineItem) BindBottleId(name string) {

	id, exists := i.BottleMap[name]
	if exists {
		i.BottleId = id
	} else {
		i.BottleId = -1
	}

}

func NewSupplierOrderLineControl(bottleNames []string, bottleMap map[string]int) *SupplierOrderLineItem {

	i := &SupplierOrderLineItem{
		SelectBottle:  widget.NewSelect(bottleNames, func(s string) {}),
		EntryQuantity: widget.NewEntry(),
		BottleMap:     bottleMap,
	}

	i.SelectBottle.PlaceHolder = " "
	i.SelectBottle.OnChanged = func(name string) {
		i.BindBottleId(name)
	}

	return i
}
