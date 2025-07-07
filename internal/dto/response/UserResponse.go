package response

type UserResponse struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Email    string  `json:"email"`
	Profile  ProfileResponse
	Roles    []RoleResponse `json:"roles"`
}

type UserDetailResponse struct {
	ID             string  `json:"id"`
	Username       *string `json:"username"`
	Email          string  `json:"email"`
	GoogleID       *string `json:"google_id"`
	EmailVerified  bool    `json:"email_verified"`
	CreatedByAdmin bool    `json:"created_by_admin"`
	Profile        ProfileDetailResponse
	Roles          []RoleResponse `json:"roles"`
}
