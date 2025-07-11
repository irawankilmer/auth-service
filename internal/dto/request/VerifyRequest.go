package request

type VerifyRegisterByAdminRequest struct {
	Username        string `json:"username" binding:"required,excludesall= "`
	Password        string `json:"password" binding:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

func (v *VerifyRegisterByAdminRequest) Sanitize() map[string]any {
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
