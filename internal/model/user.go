package model

// CreateUser struct represents mandatory user information for registration
// swagger:model
type CreateUser struct {
	// The Name of a user
	// example: Some name
	// required: true
	UserName string `json:"userName"`
	// The email of a user
	// example: qwerty@gmail.com
	// required: true
	Email string `json:"email"`
	// The password of a user
	// example: ZAQ!2wsx
	// required: true
	Password string `json:"password"`
}

// AuthUser struct represents mandatory user information for authorisation
// swagger:model
type AuthUser struct {
	// The email of a user
	// example: qwerty@gmail.com
	// required: true
	Email string `json:"email"`
	// The password of a user
	// example: ZAQ!2wsx
	// required: true
	Password string `json:"password"`
}

// Tokens struct represents a couple of token
// swagger:model
type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshToken struct represents a  refresh token
// swagger:model
type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}
