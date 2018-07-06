package contract

import (
	"resulguldibi/color-api/entity"
)

type GetColorResponse struct {
	MixedColor      *entity.Color   `json:"mixedColor"`
	RandomColors    []*entity.Color `json:"randomColors"`
	Code            string          `json:"code"`
	RaundStartPoint int             `json:"raundStartPoint"`
	TotalPoint      int             `json:"totalPoint"`
}

type GetLevelResponse struct {
	LevelCount   int `json:"levelCount"`
	DefaultLevel int `json:"defaultLevel"`
}

type GetColorHelpResponse struct {
	SelectedColors []*entity.Color `json:"selectedColors"`
	Point          int             `json:"point"`
}

type GetColorStepHelpResponse struct {
	Color *entity.Color `json:"color"`
	Point int           `json:"point"`
}

type GetColorNameResponse struct {
	Name *ColorNameItem `json:"name"`
}

type ColorNameItem struct {
	Value           string `json:"value"`
	ClosestNamedHex string `json:"closest_named_hex"`
}

type ValidateColorsResponse struct {
	IsValid    bool `json:"isValid"`
	RaundPoint int  `json:"raundPoint"`
	TotalPoint int  `json:"totalPoint"`
}

type ValidateColorsRequest struct {
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

type GetUserRaundHistoryResponse struct {
	Attempts []*entity.UserRaundColorValidationAttempt `json:"attempts"`
}

type GetRankingRequest struct {
}
