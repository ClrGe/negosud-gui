package config

import (
	"time"
)

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

type PartialProducer struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_By"`
}

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

// TODO : wait for the Oder API to be implemented
// Define the order struct and associate json fields
type Order struct {
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

var orders []Order
var order []Order
var Bottles []Bottle
var Individual Producer
var ProducerData []PartialProducer
var Producers []Producer
