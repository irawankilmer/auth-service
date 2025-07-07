package request

type ProfileCreateRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Address  *string `json:"address"`
	Gender   *string `json:"gender"`
	Image    *string `json:"image"`
}
