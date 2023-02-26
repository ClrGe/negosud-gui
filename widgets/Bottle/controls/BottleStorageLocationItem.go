package controls

import (
	"fyne.io/fyne/v2/widget"
)

type BottleStorageLocationItem struct {
	BottleId              int
	StorageLocationId     int
	SelectStorageLocation *widget.Select
	EntryQuantity         *widget.Entry
	StorageLocationMap    map[string]int
}

func (i *BottleStorageLocationItem) Bind(bottleId int, storageLocationId int) {
	i.BottleId = bottleId
	i.StorageLocationId = storageLocationId
}

func (i *BottleStorageLocationItem) BindStorageLocationId(name string) {

	id, exists := i.StorageLocationMap[name]
	if exists {
		i.StorageLocationId = id
	} else {
		i.StorageLocationId = -1
	}

}

func NewBottleStorageLocationControl(bottleNames []string, storageLocationMap map[string]int) *BottleStorageLocationItem {

	i := &BottleStorageLocationItem{
		SelectStorageLocation: widget.NewSelect(bottleNames, func(s string) {}),
		EntryQuantity:         widget.NewEntry(),
		StorageLocationMap:    storageLocationMap,
	}

	i.SelectStorageLocation.PlaceHolder = " "
	i.SelectStorageLocation.OnChanged = func(name string) {
		i.BindStorageLocationId(name)
	}

	return i
}
