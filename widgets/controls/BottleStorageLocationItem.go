package controls

import (
	"fyne.io/fyne/v2/widget"
)

type BottleStorageLocationItem struct {
	BottleId          int
	StorageLocationId int
	SelectEntryBottle *widget.SelectEntry
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
		SelectEntryBottle: widget.NewSelectEntry(bottleNames),
		EntryQuantity:     widget.NewEntry(),
		BottleMap:         bottleMap,
	}

	i.SelectEntryBottle.OnChanged = func(name string) {
		i.BindBottleId(name)
	}

	return i
}
