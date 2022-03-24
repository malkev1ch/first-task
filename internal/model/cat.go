// Package model represent objects structure in application
package model

import "time"

type Cat struct {
	ID         *string    `json:"id" bson:"_id"`
	Name       *string    `json:"name" bson:"name"`
	DateBirth  *time.Time `json:"dateBirth" bson:"dateBirth"`
	Vaccinated *bool      `json:"vaccinated" bson:"vaccinated"`
	ImagePath  *string    `json:"imagePath,omitempty" bson:"imagePath"`
}
