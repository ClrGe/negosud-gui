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
	ID    string
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"details"`
	Role  string `json:"created_By"`
}

// ############################################
// ############### PRODUCERS ##################
// ############################################

var Individual Producer
var ProducerData []PartialProducer
var Producers []Producer

// Producer struct holds information about a producer
type Producer struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Details   string      `json:"details"`
	CreatedAt interface{} `json:"created_At"`
	UpdatedAt time.Time   `json:"updated_At"`
	CreatedBy string      `json:"created_By"`
	UpdatedBy string      `json:"updated_By"`
	Bottles   interface{} `json:"Bottles"`
	Region    interface{} `json:"region"`
}

// PartialProducer holds only the necessary data for the table (= needs string only)
type PartialProducer struct {
	Id        int `json:"id"`
	ID        string
	Name      string `json:"name"`
	CreatedBy string `json:"created_By"`
}

// ############################################
// ################# BOTTLES ##################
// ############################################

var Bottles []Bottle
var BottleData []PartialBottle
var IndBottle Bottle

// Bottle struct holds information about a bottle of wine.
type Bottle struct {
	ID                int         `json:"id"`
	FullName          string      `json:"full_Name"`
	Description       string      `json:"description"`
	Label             string      `json:"label"`
	Volume            int         `json:"volume"`
	Picture           string      `json:"picture"`
	YearProduced      int         `json:"year_Produced"`
	AlcoholPercentage int         `json:"alcohol_Percentage"`
	CurrentPrice      int         `json:"current_Price"`
	CreatedAt         time.Time   `json:"created_At"`
	UpdatedAt         time.Time   `json:"updated_At"`
	CreatedBy         string      `json:"created_By"`
	UpdatedBy         string      `json:"updated_By"`
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
	FullName          string `json:"full_Name"`
	Description       string `json:"description"`
	Label             string `json:"label"`
	VolumeInt         int    `json:"volume"`
	Picture           string `json:"picture"`
	YearProduced      int    `json:"year_Produced"`
	AlcoholPercentage int    `json:"alcohol_Percentage"`
	CurrentPrice      int    `json:"current_Price"`
}

// ############################################
// ################## ORDERS ##################
// ############################################

var Orders []Order
var order []Order
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
