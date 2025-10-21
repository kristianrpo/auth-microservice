package request

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	IDCitizen int    `json:"id_citizen" validate:"required,gt=0"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Name      string `json:"name" validate:"required,min=2"`
}
