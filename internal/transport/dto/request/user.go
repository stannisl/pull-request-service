package request

type UserIsActive struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

