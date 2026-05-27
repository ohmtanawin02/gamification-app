package dto

import (
	"github.com/go-playground/validator/v10"

	"gamification-app/internal/users/domain"
)

type LoginRequest struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=100"`
}

func (r *LoginRequest) Validate(v *validator.Validate) error {
	return v.Struct(r)
}

func (r *LoginRequest) ToDomainInput() domain.LoginInput {
	return domain.LoginInput{Nickname: r.Nickname}
}
