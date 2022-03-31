// Package model represent objects structure in application
package model

import "time"

// swagger:model Cat
type Cat struct {
	// The UUID of a cat
	// required: false
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	ID string `json:"id" bson:"_id"`
	// The Name of a cat
	// example: Some name
	// required: true
	Name string `json:"name" bson:"name"`
	// The birthdate of a cat
	// example: 2018-09-22T12:42:31Z
	// required: true
	DateBirth time.Time `json:"dateBirth" bson:"dateBirth"`
	// The status of vaccination of a cat
	// example: true
	// required: true
	Vaccinated bool `json:"vaccinated" bson:"vaccinated"`
	// The image path of a cat
	// example: 1c219a3f-a959-4395-81f0-4e735040ed61.webp
	ImagePath string `json:"imagePath,omitempty" bson:"imagePath"`
}

// CreateCat is the struct for adding a cat
// swagger:model
type CreateCat struct {
	// The Name of a cat
	// example: Some name
	// required: true
	Name string `json:"name" bson:"name" validate:"required"`
	// The birthdate of a cat
	// example: 2018-09-22T12:42:31Z
	// required: true
	DateBirth time.Time `json:"dateBirth" bson:"dateBirth" validate:"required"`
	// The status of vaccination of a cat
	// example: true
	// required: true
	Vaccinated bool `json:"vaccinated" bson:"vaccinated"`
}

// UpdateCat is the struct for update a cat
// swagger:model
type UpdateCat struct {
	// The Name of a cat
	// example: Some name
	// required: true
	Name *string `json:"name" bson:"name" validate:""`
	// The birthdate of a cat
	// example: 2018-09-22T12:42:31Z
	// required: true
	DateBirth *time.Time `json:"dateBirth" bson:"dateBirth" validate:""`
	// The status of vaccination of a cat
	// example: true
	// required: true
	Vaccinated *bool `json:"vaccinated" bson:"vaccinated" validate:""`
}
