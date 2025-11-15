package dto

import "github.com/stannisl/avito-test/internal/domain"

type User struct {
	Id       domain.UserID `json:"user_id" binding:"required"`
	Username string        `json:"username" binding:"required"`
	IsActive bool          `json:"is_active" binding:"required"`
}

func (u *User) FromModel(model *domain.User) *User {
	u.Id = model.Id
	u.Username = model.Username
	u.IsActive = model.IsActive
	return u
}

func (u *User) ToModel() domain.User {
	return domain.User{
		Id:       u.Id,
		Username: u.Username,
		IsActive: u.IsActive,
	}

}
