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

type GetGoogleOAuthTokenRequest struct {
	Token string `json:"token"`
}

type GetGoogleOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type GetRankingResponse struct {
	RaundPoint int `json:"raundPoint"`
	TotalPoint int `json:"totalPoint"`
}

type GetRankingRequest struct {
}
