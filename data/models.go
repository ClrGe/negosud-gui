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
	Id       int    `json:"id"`
	Name     string `json:"firstName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	//Role     string `json:"role"`
}

// ############################################
// ############### PRODUCERS ##################
// ############################################

var Individual Producer
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
	ID                int         `json:"id"`
	FullName          string      `json:"fullName"`
	Description       string      `json:"description"`
	WineType          string      `json:"wineType"`
	Volume            int         `json:"volume"`
	Picture           string      `json:"picture"`
	YearProduced      int         `json:"yearProduced"`
	AlcoholPercentage int         `json:"alcoholPercentage"`
	CurrentPrice      int         `json:"currentPrice"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
	CreatedBy         string      `json:"createdBy"`
	UpdatedBy         string      `json:"updatedBy"`
	BottleLocations   interface{} `json:"bottleLocations"`
	BottleGrapes      interface{} `json:"bottleGrapes"`
	Producer          interface{} `json:"producer"`
}

// PartialBottle holds only the necessary data for the table (= needs string only)
type PartialBottle struct {
	Id                int `json:"id"`
	ID                string
	Price             string
	Volume            string
	Alcohol           string
	Year              string
	FullName          string `json:"fullName"`
	Description       string `json:"description"`
	WineType          string `json:"wineType"`
	VolumeInt         int    `json:"volume"`
	Picture           string `json:"picture"`
	YearProduced      int    `json:"yearProduced"`
	AlcoholPercentage int    `json:"alcoholPercentage"`
	CurrentPrice      int    `json:"currentPrice"`
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
