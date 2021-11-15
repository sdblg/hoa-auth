package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique"`
	Addresses []Address `json:"addresses,omitempty"`
	Role      string    `json:"role"`
	Password  []byte    `json:"-"`
}

type Address struct {
	gorm.Model
	UserID     uint
	Street     string `json:"street"`
	Apartment  string `json:"apartment"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}
