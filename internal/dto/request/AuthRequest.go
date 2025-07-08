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
