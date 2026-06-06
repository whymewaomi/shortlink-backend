package core_dto

type RegisterUserDto struct {
	Username        string `json:"username" validate:"required,min=4,max=32"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=48"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=6,max=48"`
}

type UserDtoResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUserDto struct {
	Username string
	Password string
}
type ShortLinkDto struct {
	Link string `json:"link" validate:"required"`
}
