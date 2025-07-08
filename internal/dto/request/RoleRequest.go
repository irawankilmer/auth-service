package request

type RoleRequest struct {
	Roles []string `json:"roles" binding:"required"`
}

func (r *RoleRequest) Sanitize() map[string]any {
	return map[string]any{
		"roles": r.Roles,
	}
}
