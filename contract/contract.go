package contract

import (
	"resulguldibi/color-api/entity"
)

type GetColorResponse struct {
	MixedColor   *entity.Color   `json:"mixedColor"`
	RandomColors []*entity.Color `json:"randomColors"`
	Code         string          `json:"code"`
}

type ValidateColorsResponse struct {
	IsValid bool `json:"isValid"`
}

type ValidateColorsRequest struct {
	MixedColor     *entity.Color   `json:"mixedColor"`
	SelectedColors []*entity.Color `json:"selectedColors"`
}

type GetRankingResponse struct {
	RaundPoint int `json:"raundPoint"`
	TotalPoint int `json:"totalPoint"`
}

type GetRankingRequest struct {
}
