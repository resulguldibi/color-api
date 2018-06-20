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

	return response, err
}

func (service ColorService) ValidateColors(colors []*entity.Color, color *entity.Color) (*contract.ValidateColorsResponse, error) {
	response := &contract.ValidateColorsResponse{}
	var err error

	mixedColor := util.GenerateMixColor(colors)
	response.IsValid = util.IsColorsEquals(mixedColor, color)

	return response, err
}
