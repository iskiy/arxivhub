package models

type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=5,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewRegisterResponse(user User) RegisterUserResponse {
	return RegisterUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

type UpdateUserEmailParams struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
