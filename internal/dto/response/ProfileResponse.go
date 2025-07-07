package response

type ProfileResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type ProfileDetailResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Address  *string `json:"address"`
	Gender   *string `json:"gender"`
	Image    *string `json:"image"`
}
