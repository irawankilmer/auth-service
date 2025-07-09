package request

type UserCreateRequest struct {
	Email string   `json:"email" binding:"required,email"`
	Roles []string `json:"roles" binding:"required"`
}

func (u *UserCreateRequest) Sanitize() map[string]any {
	return map[string]any{
		"email": u.Email,
		"roles": u.Roles,
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

func (u *UserUpdateEmailRequest) Sanitize() map[string]any {
	return map[string]any{
		"email": u.Email,
	}
}
