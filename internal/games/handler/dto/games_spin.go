package dto

import (
	"gamification-app/internal/games/domain"

	"github.com/go-playground/validator/v10"
)

type SpinRequest struct {
	PointsEarned int `json:"points_earned" validate:"required,min=1"`
}

func (r *SpinRequest) Validate(v *validator.Validate) error {
	return v.Struct(r)
}

func (r *SpinRequest) ToDomainInput(userID uint) domain.SpinInput {
	return domain.SpinInput{
		UserID:       userID,
		PointsEarned: r.PointsEarned,
	}
}

type SpinResponse struct {
	TotalPoints int `json:"total_points"`
}

func ToSpinResponse(result *domain.SpinResult) SpinResponse {
	return SpinResponse{TotalPoints: result.TotalPoints}
}
