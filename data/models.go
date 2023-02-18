package data

import (
	"time"
)

// ############################################
// ################## USERS ###################
// ############################################

var Users []User

// User struct holds information about a user
type User struct {
	ID        string `json:"-"`
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	//Role     string `json:"role"`
}

// ############################################
// ############### PRODUCERS ##################
// ############################################

var IndProducer Producer
var ProducerData []PartialProducer

// Producer struct holds information about a producer
type Producer struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Details   string      `json:"details"`
	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
	Bottles   interface{} `json:"Bottles"`
	Region    interface{} `json:"region"`
}

// PartialProducer holds only the necessary data for the table (= needs string only)
type PartialProducer struct {
	Id        int `json:"id"`
	ID        string
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
}

// ############################################
// ################# BOTTLES ##################
// ############################################

var BottleData []PartialBottle
var IndBottle Bottle
var Bottles []Bottle

// Bottle struct holds information about a bottle of wine.
type Bottle struct {
	ID                     int                     `json:"id"`
	FullName               string                  `json:"fullName"`
	Description            string                  `json:"description"`
	WineType               string                  `json:"wineType"`
	Volume                 int                     `json:"volume"`
	Picture                string                  `json:"picture"`
	YearProduced           int                     `json:"yearProduced"`
	AlcoholPercentage      float32                 `json:"alcoholPercentage"`
	CurrentPrice           float32                 `json:"currentPrice"`
	CreatedAt              time.Time               `json:"createdAt"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	CreatedBy              string                  `json:"createdBy"`
	UpdatedBy              string                  `json:"updatedBy"`
	BottleStorageLocations []BottleStorageLocation `json:"bottleStorageLocations"`
	BottleGrapes           interface{}             `json:"bottleGrapes"`
	BottleSuppliers        interface{}             `json:"bottleSuppliers"`
	Producer               interface{}             `json:"producer"`
}

// PartialBottle holds only the necessary data for the table (= needs string only)
type PartialBottle struct {
	Id                int `json:"id"`
	ID                string
	Price             string
	Volume            string
	Alcohol           string
	Year              string
	FullName          string  `json:"fullName"`
	Description       string  `json:"description"`
	WineType          string  `json:"wineType"`
	VolumeInt         int     `json:"volume"`
	Picture           string  `json:"picture"`
	YearProduced      int     `json:"yearProduced"`
	AlcoholPercentage float32 `json:"alcoholPercentage"`
	CurrentPrice      float32 `json:"currentPrice"`
}

// ############################################
// ################## ORDERS ##################
// ############################################

var Orders []Order
var CustomerOrders []CustomerOrder

type Order struct {
	ID          string
	Product     string
	Quantity    string
	Producer    string
	Date        string
	Comment     string
	Status      string
	Id          int `json:"id"`
	ProductId   int `json:"bottle_id"`
	QuantityInt int `json:"quantity"`
	ProducerId  int `json:"producer_id"`
}

type CustomerOrder struct {
	ID       string
	Client   string
	Product  string
	Quantity string
	Producer string
	Date     string
	Status   string
}

// ############################################
// ################## STORAGE LOCATIONS ###################
// ############################################

var StorageLocations []StorageLocation
var IndStorageLocation StorageLocation
var StorageLocationsData []PartialStorageLocation

// StorageLocation struct holds information about a user
type StorageLocation struct {
	ID                     int                     `json:"id"`
	Name                   string                  `json:"name"`
	CreatedAt              interface{}             `json:"createdAt"`
	UpdatedAt              interface{}             `json:"updatedAt"`
	CreatedBy              string                  `json:"createdBy"`
	UpdatedBy              string                  `json:"updatedBy"`
	BottleStorageLocations []BottleStorageLocation `json:"bottleStorageLocations"`
}

// PartialStorageLocation holds only the necessary data for the table (= needs string only)
type PartialStorageLocation struct {
	Id        int `json:"id"`
	ID        string
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

// ############################################
// ################## BOTTLE STORAGE LOCATIONS ###################
// ############################################

var BottleStorageLocations []BottleStorageLocation
var IndBottleStorageLocation BottleStorageLocation
var BottleStorageLocationData []PartialBottleStorageLocation

// BottleStorageLocation struct holds information about a user
type BottleStorageLocation struct {
	ID              int             `json:"id"`
	Bottle          Bottle          `json:"Bottle"`
	StorageLocation StorageLocation `json:"StorageLocation"`
	Quantity        int             `json:"Quantity"`
	CreatedAt       interface{}     `json:"createdAt"`
	UpdatedAt       interface{}     `json:"updatedAt"`
	CreatedBy       string          `json:"createdBy"`
	UpdatedBy       string          `json:"updatedBy"`
}

// PartialBottleStorageLocation holds only the necessary data for the table (= needs string only)
type PartialBottleStorageLocation struct {
	Id                  int `json:"id"`
	ID                  string
	BottleName          string
	StorageLocationName string
	Quantity            int    `json:"Quantity"`
	Name                string `json:"name"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
	CreatedBy           string `json:"createdBy"`
	UpdatedBy           string `json:"updatedBy"`
}
