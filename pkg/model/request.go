package model

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,min=13"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	DeviceInfo string `json:"device_info"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceInfo   string `json:"device_info"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	LogoutAll    bool   `json:"logout_all"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"required,min=13"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
