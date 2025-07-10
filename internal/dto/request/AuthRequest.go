package request

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (l *LoginRequest) Sanitize() map[string]any {
	return map[string]any{
		"identifier": l.Identifier,
	}
}

type RegisterRequest struct {
	FullName        string   `json:"full_name" binding:"required"`
	Username        string   `json:"username" binding:"required,excludesall= "`
	Email           string   `json:"email" binding:"required,email"`
	Password        string   `json:"password" binding:"required,min=6"`
	ConfirmPassword string   `json:"confirm_password" binding:"required"`
	Roles           []string `json:"roles" binding:"required"`
}

func (r *RegisterRequest) Sanitize() map[string]any {
	return map[string]any{
		"full_name": r.FullName,
		"username":  r.Username,
		"email":     r.Email,
	}
}
