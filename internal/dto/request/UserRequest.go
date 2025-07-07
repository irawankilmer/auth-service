package request

type UserCreateRequest struct {
	Username string `json:"username" binding:"required,excludesall= "`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"min=6"`
	Profile  ProfileCreateRequest
	Roles    []string `json:"roles" binding:"required"`
}

func (u *UserCreateRequest) Sanitize() map[string]any {
	return map[string]any{
		"username": u.Username,
		"email":    u.Email,
	}
}

type UserUpdateUsernameRequest struct {
	Username string `json:"username" binding:"required,excludesall= "`
}

func (u *UserUpdateUsernameRequest) Sanitize() map[string]any {
	return map[string]any{
		"username": u.Username,
	}
}

type UserUpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
