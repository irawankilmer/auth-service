package model

type ProfileModel struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	FullName string  `json:"full_name"`
	Address  *string `json:"address"`
	Gender   *string `json:"gender"`
	Image    *string `json:"image"`
}
