package controls

import (
	"fyne.io/fyne/v2/widget"
)

type BottleStorageLocationItem struct {
	BottleId          int
	StorageLocationId int
	SelectBottle      *widget.Select
	EntryQuantity     *widget.Entry
	BottleMap         map[string]int
}

func (i *BottleStorageLocationItem) Bind(bottleId int, storageLocationId int) {
	i.BottleId = bottleId
	i.StorageLocationId = storageLocationId
}

func (i *BottleStorageLocationItem) BindBottleId(name string) {

	id, exists := i.BottleMap[name]
	if exists {
		i.BottleId = id
	} else {
		i.BottleId = -1
	}

}

func NewBottleStorageLocationControl(bottleNames []string, bottleMap map[string]int) *BottleStorageLocationItem {

	i := &BottleStorageLocationItem{
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
