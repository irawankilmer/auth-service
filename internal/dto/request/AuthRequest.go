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
	Username string `json:"username" binding:"required,excludesall= "`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"min=6"`
	Profile  ProfileCreateRequest
	Roles    []string `json:"roles" binding:"required"`
}
