package dto

import (
	"gamification-app/internal/rewards/domain"

	"github.com/go-playground/validator/v10"
)

type ClaimRewardRequest struct {
	RewardID uint `json:"reward_id" validate:"required,min=1"`
}

func (r *ClaimRewardRequest) Validate(v *validator.Validate) error {
	return v.Struct(r)
}

func (r *ClaimRewardRequest) ToDomainInput(userID uint) domain.ClaimRewardInput {
	return domain.ClaimRewardInput{
		UserID:   userID,
		RewardID: r.RewardID,
	}
}
