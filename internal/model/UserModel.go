package model

type UserModel struct {
	ID             string
	Username       *string
	Email          string
	Password       *string
	TokenVersion   string
	EmailVerified  bool
	CreatedByAdmin bool
	GoogleID       *string
	Profile        ProfileModel
	Roles          []RoleModel
}
