// Package model represent objects structure in application
package model

import "time"

type Cat struct {
	Name       *string    `json:"name"`
	DateBirth  *time.Time `json:"dateBirth"`
	Vaccinated *bool      `json:"vaccinated"`
	ImagePath  *string    `json:"imagePath,omitempty"`
}
