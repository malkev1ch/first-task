package domain

import "time"

type CreateCat struct {
	Name       string    `json:"name"`
	DateBirth  time.Time `json:"dateBirth"`
	Vaccinated bool      `json:"vaccinated"`
}

type UpdateCat struct {
	Name       *string    `json:"name"`
	DateBirth  *time.Time `json:"dateBirth"`
	Vaccinated *bool      `json:"vaccinated"`
}

type Cat struct {
	Name       string    `json:"name"`
	DateBirth  time.Time `json:"dateBirth"`
	Vaccinated bool      `json:"vaccinated"`
}
