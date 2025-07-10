package request

type VerifyRegisterRequest struct {
	Username        string `json:"username" binding:"required,excludesall= "`
	Password        string `json:"password" binding:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

func (v *VerifyRegisterRequest) Sanitize() map[string]any {
	return map[string]any{
		"username": v.Username,
	}
}

type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

func (v *VerifyRequest) Sanitize() map[string]any {
	return map[string]any{
		"token": v.Token,
	}
}
