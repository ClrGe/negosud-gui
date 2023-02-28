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
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Details string      `json:"details"`
	Address *Address    `json:"address"`
	Region  interface{} `json:"region"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
	Bottles   interface{} `json:"Bottles"`
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
	ThresholdToOrder       int                     `json:"thresholdToOrder"`
	QuantityMinimumToOrder int                     `json:"quantityMinimumToOrder"`
	CurrentPrice           float32                 `json:"currentPrice"`
	CreatedAt              interface{}             `json:"createdAt"`
	UpdatedAt              interface{}             `json:"updatedAt"`
	CreatedBy              string                  `json:"createdBy"`
	UpdatedBy              string                  `json:"updatedBy"`
	BottleStorageLocations []BottleStorageLocation `json:"bottleStorageLocations"`
	BottleGrapes           interface{}             `json:"bottleGrapes"`
	BottleSuppliers        interface{}             `json:"bottleSuppliers"`
	Producer               interface{}             `json:"producer"`
	CustomerPrice          float32                 `json:"customerPrice"`
	SupplierPrice          float32                 `json:"supplierPrice"`
	IsSelected             bool
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
	CustomerPrice     float32 `json:"customerPrice"`
	SupplierPrice     float32 `json:"supplierPrice"`
}

// ############################################
// ################## CUSTOMERORDERS ##################
// ############################################

var CustomerOrders []CustomerOrder
var CustomerOrderData []PartialCustomerOrder

type CustomerOrder struct {
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

type PartialCustomerOrder struct {
	ID       string
	Client   string
	Product  string
	Quantity string
	Producer string
	Date     string
	Status   string
}

//############################################
//################## SUPPLIER ORDERS ##################
//############################################

var SupplierOrders []SupplierOrder
var IndSupplierOrder SupplierOrder
var SupplierOrderData []PartialSupplierOrder

type SupplierOrder struct {
	ID             int                 `json:"id"`
	Reference      string              `json:"reference"`
	DateOrder      interface{}         `json:"dateOrder"`
	DateDelivery   interface{}         `json:"dateDelivery"`
	DeliveryStatus int                 `json:"deliveryStatus"`
	Supplier       *Supplier           `json:"supplier"`
	Lines          []SupplierOrderLine `json:"Lines"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

type PartialSupplierOrder struct {
	Id        int `json:"id"`
	ID        string
	Reference string `json:"reference"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

//############################################
//################## SUPPLIER ORDER LINES ##################
//############################################

var SupplierOrderLines []CustomerOrder
var SupplierOrderLineData []PartialCustomerOrder

type SupplierOrderLine struct {
	ID            int           `json:"id"`
	BottleId      int           `json:"BottleId"`
	Bottle        Bottle        `json:"Bottle"`
	SupplierOrder SupplierOrder `json:"SupplierOrder"`
	Quantity      int           `json:"Quantity"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

type PartialSupplierOrderLine struct {
	Id                     int `json:"id"`
	ID                     string
	BottleName             string
	SupplierOrderReference string

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

// ############################################
// ################## STORAGE LOCATIONS ###################
// ############################################

var StorageLocations []StorageLocation
var IndStorageLocation StorageLocation
var StorageLocationData []PartialStorageLocation

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

// ############################################
// ################## SUPPLIER ###################
// ############################################

var Suppliers []Supplier
var IndSupplier Supplier
var SupplierData []PartialSupplier

// Supplier struct holds information about a Supplier
type Supplier struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Details   string      `json:"details"`
	Email     string      `json:"email"`
	Address   *Address    `json:"address"`
	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

// PartialSupplier holds only the necessary data for the table (= needs string only)
type PartialSupplier struct {
	Id        int `json:"id"`
	ID        string
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

// ############################################
// ################## ADDRESS ###################
// ############################################

type Address struct {
	ID           int    `json:"id"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	CityId       int    `json:"cityId"`
	City         *City  `json:"city"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

// ############################################
// ################## CITY ###################
// ############################################

var Cities []City
var IndCity City
var CityData []PartialCity

type City struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	ZipCode int      `json:"zipCode"`
	Country *Country `json:"country"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

type PartialCity struct {
	Id      int `json:"id"`
	ID      string
	Name    string `json:"name"`
	ZipCode int    `json:"zipCode"`
}

// ############################################
// ################## Country ###################
// ############################################

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

// ############################################
// ################## WINELABEL ###################
// ############################################

var IndWineLabel WineLabel
var WineLabelData []PartialWineLabel

type WineLabel struct {
	ID    int    `json:"id"`
	Label string `json:"label"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

type PartialWineLabel struct {
	Id    int `json:"id"`
	ID    string
	Label string `json:"label"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}

// ############################################
// ################## GRAPE ###################
// ############################################

var IndGrape Grape
var GrapeData []PartialGrape

type Grape struct {
	ID        int    `json:"id"`
	GrapeType string `json:"grapeType"`

	CreatedAt interface{} `json:"createdAt"`
	UpdatedAt interface{} `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
	UpdatedBy string      `json:"updatedBy"`
}

type PartialGrape struct {
	Id        int `json:"id"`
	ID        string
	GrapeType string `json:"grapeType"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedBy string `json:"updatedBy"`
}
