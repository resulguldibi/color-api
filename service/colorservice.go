package service

import (
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/util"
)

func (service ColorService) GetRandomColors(level int64) (*contract.GetColorResponse, error) {
	response := &contract.GetColorResponse{}
	var err error
	var count int64 = 5*level + 1

	randomColors := make([]*entity.Color, 0, count)

	for i := int64(0); i < level; i++ {
		randomColor := util.GenerateRandomColor()
		if !util.IsColorExist(randomColors, randomColor) {
			randomColor.IsSelected = true
			randomColors = append(randomColors, randomColor)
		}
	}

	mixedColor := util.GenerateMixColor(randomColors)

	mixColorList := make([]*entity.Color, 0, 1)
	mixColorList = append(mixColorList, mixedColor)

	for i := int64(0); i < 4*level; i++ {
		randomColor := util.GenerateRandomColor()
		if !util.IsColorExist(randomColors, randomColor) && !util.IsColorExist(mixColorList, randomColor) {
			randomColors = append(randomColors, randomColor)
		}
	}

	finalRandomColors := make([]*entity.Color, 0, len(randomColors))
	for _, randomColor := range randomColors {
		finalRandomColors = append(finalRandomColors, randomColor)
	}

	for i := int64(0); i < level; i++ {

		random := i
		for random == i {
			random = int64(util.GenerateRandomNumber(int(5 * level)))
		}

		tmp := finalRandomColors[i]
		finalRandomColors[i] = finalRandomColors[random]
		finalRandomColors[random] = tmp
	}

	response.MixedColor = mixedColor
	response.RandomColors = finalRandomColors
	response.Code = util.GenerateGuid()

	//1. get user-id from client in http request header with key "Authorization" in JWT format.
	//2. update user point calculation key with this guid (store this data in redis -> hmset user-point-calculation user-id "1234" guid "12s12-12sas-3asw12-12sa1")

	return response, err
}

func (service ColorService) ValidateColors(colors []*entity.Color, color *entity.Color) (*contract.ValidateColorsResponse, error) {
	response := &contract.ValidateColorsResponse{}
	var err error

	//1. get user-id from client in http request header with key "Authorization" in JWT format.
	//2. get code from client (code was sended to client in /colors response) in /validate request to calculate user point (client should send this guid in /validate request)
	//3. increment user step number in every /validate request (store this data in redis        -> hmset user-step-number user-id "1234" guid "12s12-12sas-3asw12-12sa1" new-step-number "10")

	isMatched := util.IsMatchedColors(colors, color)

	if !isMatched {
		//3.1.a if step number is reached to max retry number, then game is over.
	} else {
		//3.2.a calculate point with generated point algorithm.
		//3.2.b update user point if /validate endpoint return success code(store this data in redis   -> hmset user-point user-id "1234" guid "12s12-12sas-3asw12-12sa1" point "100")
	}

	response.IsValid = isMatched
	return response, err
}
