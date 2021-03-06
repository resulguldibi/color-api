package service

import (
	"encoding/json"
	"fmt"
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/types"
	"resulguldibi/color-api/util"
	"strconv"
)

/*

hmget user-raund-key 12345
hmget user-raund-level 12345
hmget user-total-point 12345
hmget user-raund-point 12345
hmget user-raund-point "56dd5068-20ce-4f6d-845b-ea4990008bac"
hmget user-raund-step-number 56dd5068-20ce-4f6d-845b-ea4990008bac
hmget user-raund-generated-random-colors 56dd5068-20ce-4f6d-845b-ea4990008bac


*/

const raundStartPoint int = 20
const defaultLevelNumber = 2
const totalLevelCount = 4
const defaultRaundCountPerStage = 10

func (service *ColorService) GetUserStageInfo(userId string) (*contract.GetUserStageInfoResponse, error) {
	response := &contract.GetUserStageInfoResponse{}

	stages := make([]*entity.StageInfo, 0, totalLevelCount)

	for stage := 1; stage <= totalLevelCount; stage++ {
		stageCount, err := service.getUserStageInfo(userId, stage)

		if err != nil {
			return nil, err
		}
		stages = append(stages, &entity.StageInfo{Stage: stage, UserStageCount: stageCount, DefaultRaundCountPerStage: defaultRaundCountPerStage})
	}

	response.Stages = stages

	return response, nil
}

func (service *ColorService) GetUserRaundHistory(userId string, sendedKey string) (*contract.GetUserRaundHistoryResponse, error) {
	response := &contract.GetUserRaundHistoryResponse{}
	var key string
	var err error

	key, err = service.getUserRaundKey(userId)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if key != sendedKey {
		return nil, types.NewBusinessException("invalid key", "exp.invalidkey")
	}

	attempts, err := service.getUserRaundColorValidationAttempts(key)

	if err != nil {
		return nil, err
	}

	response.Attempts = attempts

	return response, nil
}

func (service *ColorService) ValidateColors(userId string, sendedKey string, colors []*entity.Color) (*contract.ValidateColorsResponse, error) {
	response := &contract.ValidateColorsResponse{}
	var err error
	var key string

	var isColorsValid bool = false

	key, err = service.getUserRaundKey(userId)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if key != sendedKey {
		return nil, types.NewBusinessException("invalid key", "exp.invalidkey")
	}

	raundPoint, err := service.getUserRaundPoint(userId, key)
	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if raundPoint == 0 {
		return nil, types.NewBusinessException("not enough point", "exp.not.enough.point")
	}

	var mixedColor *entity.Color
	mixedColor, err = service.getUserRaundGeneratedMixedColor(key)

	if err != nil {
		return nil, err
	}

	allColors := append(colors, mixedColor)

	isColorsValid, err = service.validateSendedColors(userId, key, allColors)

	if err != nil {
		return nil, err
	}

	if !isColorsValid {
		return nil, types.NewBusinessException("invalid colors", "exp.invalidcolors")
	}

	attempt := &entity.UserRaundColorValidationAttempt{}
	attempt.SendedColors = colors
	attempt.MixedColor = util.GenerateMixColor(colors)

	err = service.addUserRaundColorValidationAttempts(key, attempt)

	if err != nil {
		return nil, err
	}

	err = service.checkUserRaundStepNumber(userId, key)
	if err != nil {
		return nil, err
	}

	err = service.updateUserRaundStepNumber(userId, key)
	if err != nil {
		return nil, err
	}

	isMatched := util.IsMatchedColors(colors, mixedColor)

	if !isMatched {

		err = service.checkUserRaundStepNumber(userId, key)
		if err != nil {

			switch err.(type) {
			case *types.GameOverException:
				innerErr := service.setUserRaundPoint(userId, key, 0)
				if innerErr != nil {
					return nil, innerErr
				}

				innerErr = service.setUserTotalPoint(userId, 0)
				if innerErr != nil {
					return nil, innerErr
				}
			}

			return nil, err
		}

		raundPoint = raundPoint - 1

		err = service.setUserRaundPoint(userId, key, raundPoint)

		if err != nil {
			return nil, types.NewBusinessException("system exception", "exp.systemexception")
		}

		response.RaundPoint = raundPoint

	} else {

		err = service.updateUserTotalPoint(userId, raundPoint)
		if err != nil {
			return nil, err
		}

		response.RaundPoint = raundPoint
	}

	response.IsValid = isMatched
	totalPoint, err := service.getUserTotalPoint(userId)
	if err == nil {
		response.TotalPoint = totalPoint
	}

	if isMatched {
		err = service.deleteUserRaundKey(userId)
		if err != nil {
			return nil, err
		}

		err = service.deleteExistingKeyRelatedData(key)

		if err != nil {
			return nil, err
		}

		err = service.setUserRaundValidation(key, true)

		if err != nil {
			return nil, err
		}

		var level int
		level, err = service.getUserRaundLevel(userId)

		if err != nil {
			return nil, err
		}

		err = service.incrementUserStageInfo(userId, level, 1)

		if err != nil {
			return nil, err
		}

	}

	return response, err
}

func (service *ColorService) GetRankings(userId string, sendedKey string) (*contract.GetRankingResponse, error) {
	response := &contract.GetRankingResponse{}
	var err error
	var key string
	var raundPoint int
	var totalPoint int

	key, err = service.getUserRaundKey(userId)

	if err != nil {
		return nil, err
	}

	if key != sendedKey {
		return nil, types.NewBusinessException("invalid key", "exp.invalidkey")
	}

	raundPoint, err = service.getUserRaundPoint(userId, key)

	if err != nil {
		return nil, err
	}

	totalPoint, err = service.getUserTotalPoint(userId)

	if err != nil {
		return nil, err
	}

	response.RaundPoint = raundPoint
	response.TotalPoint = totalPoint

	return response, err
}

func (service *ColorService) GetColorName(color *entity.Color) (*contract.GetColorNameResponse, error) {
	response := &contract.GetColorNameResponse{}
	var err error

	colorName, err := service.getColorName(color)

	if err != nil {
		fmt.Println("err -> ", err)
		return nil, types.NewBusinessException("color name exception", "exp.color.name")
	}

	if colorName == "" {
		response, err = service.GetColorNameFromColorApi(color)

		if err != nil {
			fmt.Println("err -> ", err)
			return nil, types.NewBusinessException("color name exception", "exp.color.name")
		}

		if response != nil && response.Name != nil && response.Name.Value != "" {
			colorName = response.Name.Value
			_, err = service.setColorName(color, colorName)

			if err != nil {
				fmt.Println("err -> ", err)
				return nil, types.NewBusinessException("color name exception", "exp.color.name")
			}
		}
	}

	response.Name = &contract.ColorNameItem{Value: colorName}

	return response, err
}

func (service *ColorService) GetColorNameFromColorApi(color *entity.Color) (*contract.GetColorNameResponse, error) {
	response := &contract.GetColorNameResponse{}
	var err error

	path := fmt.Sprintf("http://www.thecolorapi.com/id?rgb=rgb(%d,%d,%d)", color.R, color.G, color.B)

	err = service.httpClient.Get(path).EndStruct(response)

	if err != nil {
		fmt.Println("err -> ", err)
		return nil, types.NewBusinessException("color name exception", "exp.color.name")
	}

	return response, err
}

func (service *ColorService) GetColorStepHelp(userId string, key string, selectedColors []*entity.Color) (*contract.GetColorStepHelpResponse, error) {
	response := &contract.GetColorStepHelpResponse{}
	var err error
	var color *entity.Color
	var actualColors []*entity.Color
	var userRaundPoint int
	var isRaundAlreadyValidated = false
	var level int

	isRaundAlreadyValidated, err = service.getUserRaundValidation(key)

	if err != nil {
		return nil, err
	}

	if isRaundAlreadyValidated {
		return nil, types.NewBusinessException("raund is already validated", "exp.raund.already.validated")
	}

	helpedColors, err := service.getUserRaundStepHelp(key)

	if err != nil {
		panic(err)
	}

	//check user level for max helped colors

	level, err = service.getUserRaundLevel(userId)

	if helpedColors != nil && len(helpedColors) > 0 && level == len(helpedColors) {
		return nil, types.NewBusinessException("cannot get help anymore", "exp.cannot.stephelp.anymore")
	}

	//check user raund point to hep for new color

	userRaundPoint, err = service.getUserRaundPoint(userId, key)

	if err != nil {
		panic(err)
	}

	if userRaundPoint < raundStartPoint {
		return nil, types.NewBusinessException("not enough point", "exp.not.enough.point")
	}

	actualColors, err = service.getUserRaundGeneratedSelectedColors(key)

	if err != nil {
		panic(err)
	}

	checkForHelpedColors := helpedColors != nil && len(helpedColors) > 0

	if actualColors != nil && len(actualColors) > 0 {
		for _, actualColor := range actualColors {
			if !util.IsColorExist(selectedColors, actualColor) {

				if checkForHelpedColors {
					if !util.IsColorExist(helpedColors, actualColor) {
						color = actualColor
						helpedColors = append(helpedColors, color)
						service.setUserRaundStepHelp(key, helpedColors)
						break
					}
				} else {
					color = actualColor
					helpedColors = append(helpedColors, color)
					service.setUserRaundStepHelp(key, helpedColors)
					break
				}

			}
		}
	}

	if color != nil {
		response.Color = color
		newPoint := userRaundPoint - raundStartPoint
		err = service.setUserRaundPoint(userId, key, newPoint)
		if err != nil {
			panic(err)
		}

		response.Point = newPoint
	}

	return response, err
}

func (service *ColorService) GetColorHelp(userId string, key string) (*contract.GetColorHelpResponse, error) {
	response := &contract.GetColorHelpResponse{}
	var err error
	var selectedColors []*entity.Color

	selectedColors, err = service.getUserRaundGeneratedSelectedColors(key)

	if err != nil {
		panic(err)
	}

	//reset raund point

	err = service.setUserRaundPoint(userId, key, 0)
	if err != nil {
		panic(err)
	}

	response.SelectedColors = selectedColors
	response.Point = 0

	return response, err
}

func (service *ColorService) GetLevels() (*contract.GetLevelResponse, error) {
	response := &contract.GetLevelResponse{}
	var err error

	response.LevelCount = totalLevelCount
	response.DefaultLevel = defaultLevelNumber

	return response, err
}

func (service *ColorService) GetRandomColorsWithOutUser(level int64) (*contract.GetColorResponse, error) {

	response := &contract.GetColorResponse{}
	var err error
	var count int64 = 5*level + 1

	randomColors := make([]*entity.Color, 0, count)
	selectedColors := make([]*entity.Color, 0, level)

	for i := int64(0); i < level; i++ {
		randomColor := util.GenerateRandomColor()
		if !util.IsColorExist(randomColors, randomColor) {
			randomColor.IsSelected = true
			randomColors = append(randomColors, randomColor)
			selectedColors = append(selectedColors, randomColor)
		}
	}
	response.SelectedColors = selectedColors

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
		getColorNameResponse, err := service.GetColorName(randomColor)
		if err != nil {
			fmt.Println("err ->", err)
		}
		if getColorNameResponse != nil {
			randomColor.Name = getColorNameResponse.Name.Value
		}
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

	getColorNameResponse, err := service.GetColorName(mixedColor)

	if err != nil {
		panic(err)
	}

	if getColorNameResponse != nil {
		mixedColor.Name = getColorNameResponse.Name.Value
	}

	response.MixedColor = mixedColor
	response.RandomColors = finalRandomColors

	return response, nil
}

func (service *ColorService) GetExtendRandomColorsResponse(response *contract.GetColorResponse, userId string, level int64) (*contract.GetColorResponse, error) {

	var err error
	var raundNumber int64 = 0
	raundNumber, err = service.getUserRaundNumberByLevel(userId, level)

	if err != nil {
		panic(err)
	}

	response.RaundNumber = raundNumber

	response.RaundStartPoint = int(level) * raundStartPoint
	var totalPoint int
	totalPoint, err = service.getUserTotalPoint(userId)

	if err != nil {
		panic(err)
	}

	response.TotalPoint = totalPoint

	err = service.setUserRaundPoint(userId, response.Code, int(response.RaundStartPoint))

	if err != nil {
		panic(err)
	}

	err = service.setUserRaundGeneratedSelectedColors(response.Code, response.SelectedColors)

	if err != nil {
		panic(err)
	}

	//save generated random colors to use in /validate endpoint
	allColors := append(response.RandomColors, response.MixedColor)
	err = service.setUserRaundGeneratedRandomColors(response.Code, allColors)
	if err != nil {
		panic(err)
	}

	err = service.setUserRaundGeneratedMixedColor(response.Code, response.MixedColor)

	if err != nil {
		panic(err)
	}

	//update user point calculation key with this guid (store this data in redis -> hmset user-raund "1234" "12s12-12sas-3asw12-12sa1")

	var existingKey string
	existingKey, err = service.getUserRaundKey(userId)

	if err != nil {
		panic(err)
	}

	if existingKey != "" {
		err = service.deleteExistingKeyRelatedData(existingKey)
		if err != nil {
			panic(err)
		}
	}

	err = service.setUserRaundKey(userId, response.Code)
	if err != nil {
		panic(err)
	}

	err = service.setUserRaundLevel(userId, level)

	if err != nil {
		panic(err)
	}

	return response, err
}

func (service *ColorService) GetRandomColors(userId string, level int64) (*contract.GetColorResponse, error) {
	response := &contract.GetColorResponse{}

	var err error

	response, err = service.GetRandomColorsWithOutUser(level)

	if err != nil {
		panic(err)
	}

	code := util.GenerateGuid()

	response.Code = code

	return service.GetExtendRandomColorsResponse(response, userId, level)
}

func (service *ColorService) getUserRaundNumberByLevel(userId string, level int64) (int64, error) {

	userRaundKeyByLevel := fmt.Sprintf("%s-%d", userId, level)

	updatedRaundNumber, err := service.redisClient.HIncrBy("user-raund-number-by-level", userRaundKeyByLevel, 1)

	if err != nil {
		return updatedRaundNumber, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return updatedRaundNumber, nil
}

func (service *ColorService) addUserRaundColorValidationAttempts(key string, attempt *entity.UserRaundColorValidationAttempt) error {

	attempts, err := service.getUserRaundColorValidationAttempts(key)

	if err != nil {
		return err
	}

	if attempts == nil {
		attempts = make([]*entity.UserRaundColorValidationAttempt, 0, 0)
	}

	attempts = append(attempts, attempt)

	err = service.setUserRaundColorValidationAttempts(key, attempts)

	return err
}

func (service *ColorService) setUserRaundColorValidationAttempts(key string, attempts []*entity.UserRaundColorValidationAttempt) error {

	data := make(map[string]interface{})

	attemptBytes, err := json.Marshal(attempts)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	data[key] = string(attemptBytes)

	_, err = service.redisClient.HMSet("user-raund-color-validation-attempts", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundColorValidationAttempts(key string) ([]*entity.UserRaundColorValidationAttempt, error) {

	result, err := service.redisClient.HMGet("user-raund-color-validation-attempts", key)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	attempts := []*entity.UserRaundColorValidationAttempt{}

	if result == "" {
		return attempts, nil
	}

	err = json.Unmarshal([]byte(result), &attempts)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return attempts, err
}

func (service *ColorService) deleteUserRaundColorValidationAttempts(key string) error {

	_, err := service.redisClient.HDel("user-raund-color-validation-attempts", key)

	return err
}

func (service *ColorService) setUserRaundStepHelp(key string, colors []*entity.Color) error {

	data := make(map[string]interface{})

	colorBytes, err := json.Marshal(colors)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	data[key] = string(colorBytes)

	_, err = service.redisClient.HMSet("user-raund-step-helped-colors", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundStepHelp(key string) ([]*entity.Color, error) {

	result, err := service.redisClient.HMGet("user-raund-step-helped-colors", key)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	actualColors := []*entity.Color{}

	if result == "" {
		return actualColors, nil
	}

	err = json.Unmarshal([]byte(result), &actualColors)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return actualColors, err
}

func (service *ColorService) deleteUserRaundStepHelp(key string) error {

	_, err := service.redisClient.HDel("user-raund-step-helped-colors", key)

	return err
}

func (service *ColorService) setUserRaundKey(userId string, key string) error {
	data := make(map[string]interface{})

	data[userId] = key

	_, err := service.redisClient.HMSet("user-raund-key", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserRaundKey(userId string) error {

	_, err := service.redisClient.HDel("user-raund-key", userId)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserRaundValidation(key string) error {

	_, err := service.redisClient.HDel("user-raund-validation", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getColorName(color *entity.Color) (string, error) {

	code := fmt.Sprintf("%s-%s-%s", strconv.Itoa(color.R), strconv.Itoa(color.G), strconv.Itoa(color.B))

	response, err := service.redisClient.HMGet("color-name", code)

	if err != nil {
		return response, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return response, err
}

func (service *ColorService) setColorName(color *entity.Color, name string) (string, error) {

	data := make(map[string]interface{})

	code := fmt.Sprintf("%s-%s-%s", strconv.Itoa(color.R), strconv.Itoa(color.G), strconv.Itoa(color.B))

	data[code] = name

	response, err := service.redisClient.HMSet("color-name", data)

	if err != nil {
		return response, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return response, err
}

func (service *ColorService) getUserRaundKey(userId string) (string, error) {

	response, err := service.redisClient.HMGet("user-raund-key", userId)

	if err != nil {
		return response, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return response, err
}

func (service *ColorService) setUserRaundValidation(key string, isValidated bool) error {

	data := make(map[string]interface{})

	data[key] = isValidated
	_, err := service.redisClient.HMSet("user-raund-validation", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundValidation(key string) (bool, error) {
	var isValidated bool = false
	result, err := service.redisClient.HMGet("user-raund-validation", key)

	if err != nil {
		return false, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result != "" {
		isValidated, err = strconv.ParseBool(result)

		if err != nil {
			return false, types.NewBusinessException("system exception", "exp.systemexception")
		}
	}

	return isValidated, nil
}

func (service *ColorService) setUserRaundGeneratedSelectedColors(key string, colors []*entity.Color) error {
	data := make(map[string]interface{})

	colorBytes, err := json.Marshal(colors)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	data[key] = string(colorBytes)

	_, err = service.redisClient.HMSet("user-raund-generated-selected-colors", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundGeneratedSelectedColors(key string) ([]*entity.Color, error) {

	result, err := service.redisClient.HMGet("user-raund-generated-selected-colors", key)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	actualColors := []*entity.Color{}

	err = json.Unmarshal([]byte(result), &actualColors)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return actualColors, err
}

func (service *ColorService) setUserRaundGeneratedRandomColors(key string, colors []*entity.Color) error {
	data := make(map[string]interface{})

	colorBytes, err := json.Marshal(colors)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	data[key] = string(colorBytes)

	_, err = service.redisClient.HMSet("user-raund-generated-random-colors", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) setUserRaundGeneratedMixedColor(key string, color *entity.Color) error {
	data := make(map[string]interface{})

	colorBytes, err := json.Marshal(color)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	data[key] = string(colorBytes)

	_, err = service.redisClient.HMSet("user-raund-generated-mixed-color", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundGeneratedMixedColor(key string) (*entity.Color, error) {

	result, err := service.redisClient.HMGet("user-raund-generated-mixed-color", key)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	mixedColor := &entity.Color{}

	err = json.Unmarshal([]byte(result), mixedColor)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return mixedColor, err
}

func (service *ColorService) deleteExistingKeyRelatedData(key string) error {

	err := service.deleteUserExistingRaundGeneratedRandomColors(key)

	err = service.deleteUserExistingRaundGeneratedMixedColor(key)

	err = service.deleteUserExistingRaundStepNumber(key)

	err = service.deleteUserExistingRaundPoint(key)

	err = service.deleteUserExistingRaundGeneratedSelectedColors(key)

	err = service.deleteUserRaundValidation(key)

	err = service.deleteUserRaundStepHelp(key)

	return err

}

func (service *ColorService) deleteUserExistingRaundGeneratedRandomColors(key string) error {

	_, err := service.redisClient.HDel("user-raund-generated-random-colors", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserExistingRaundGeneratedSelectedColors(key string) error {

	_, err := service.redisClient.HDel("user-raund-generated-selected-colors", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserExistingRaundGeneratedMixedColor(key string) error {

	_, err := service.redisClient.HDel("user-raund-generated-mixed-color", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserExistingRaundStepNumber(key string) error {

	_, err := service.redisClient.HDel("user-raund-step-number", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) deleteUserExistingRaundPoint(key string) error {

	_, err := service.redisClient.HDel("user-raund-point", key)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundGeneratedRandomColors(key string) ([]*entity.Color, error) {

	result, err := service.redisClient.HMGet("user-raund-generated-random-colors", key)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	actualColors := []*entity.Color{}

	err = json.Unmarshal([]byte(result), &actualColors)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return actualColors, err
}

func (service *ColorService) validateSendedColors(userId string, key string, colors []*entity.Color) (bool, error) {
	var isValid bool = true

	level, err := service.getUserRaundLevel(userId)

	if err != nil {
		return false, err
	}

	expectedColorsCount := level + 1

	if expectedColorsCount != len(colors) {
		return false, types.NewBusinessException(fmt.Sprintf("invalid colors count -> expected :%d, actual :%d", expectedColorsCount, len(colors)), "exp.invalidcolorscount")
	}

	actualColors, err := service.getUserRaundGeneratedRandomColors(key)

	if err != nil {
		return false, err
	}

	if colors != nil && len(colors) > 0 && actualColors != nil && len(actualColors) > 0 {
		for _, color := range colors {
			if !util.IsColorExist(actualColors, color) {
				isValid = false
				break
			}
		}
	}

	return isValid, err
}

func (service *ColorService) updateUserRaundStepNumber(userId string, key string) error {

	//hmset user-raund-step-number "1234" "12s12-12sas-3asw12-12sa1" "12s12-12sas-3asw12-12sa1" "10")

	step, err := service.getUserRaundStepNumber(userId, key)

	if err != nil {
		return err
	}

	step = step + 1

	err = service.setUserRaundStepNumber(userId, key, step)

	return err
}

func (service *ColorService) setUserRaundStepNumber(userId string, key string, step int) error {

	data := make(map[string]interface{})

	data[userId] = key
	data[key] = step

	_, err := service.redisClient.HMSet("user-raund-step-number", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return err

}

func (service *ColorService) getUserRaundStepNumber(userId string, key string) (int, error) {

	result, err := service.redisClient.HMGet("user-raund-step-number", key)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	step, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return step, err

}

func (service *ColorService) checkUserRaundStepNumber(userId string, key string) error {

	var maxStep int
	var err error
	var level int

	level, err = service.getUserRaundLevel(userId)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	maxStep = level * raundStartPoint

	step, err := service.getUserRaundStepNumber(userId, key)
	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	if step >= maxStep {
		return types.NewGameOverException("max retry count reached", "exp.maxretrycountreached")
	}

	return nil
}

func (service *ColorService) getUserRaundPoint(userId string, key string) (int, error) {

	// hmget user-raund-point "12s12-12sas-3asw12-12sa1"

	result, err := service.redisClient.HMGet("user-raund-point", key)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	point, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return point, err
}

func (service *ColorService) setUserRaundPoint(userId string, key string, point int) error {
	//hmset user-raund-point "1234" "12s12-12sas-3asw12-12sa1" "12s12-12sas-3asw12-12sa1" "100"
	data := make(map[string]interface{})

	data[userId] = key
	data[key] = point

	_, err := service.redisClient.HMSet("user-raund-point", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserTotalPoint(userId string) (int, error) {

	// hmget user-raund-point "12s12-12sas-3asw12-12sa1"

	result, err := service.redisClient.HMGet("user-total-point", userId)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	point, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return point, nil
}

func (service *ColorService) setUserTotalPoint(userId string, point int) error {
	//hmset user-raund-point "1234" "100"
	data := make(map[string]interface{})

	data[userId] = point

	_, err := service.redisClient.HMSet("user-total-point", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) updateUserTotalPoint(userId string, pointToAdd int) error {

	//hmset user-raund-point "1234" "100"

	totalPoint, err := service.getUserTotalPoint(userId)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	totalPoint = totalPoint + pointToAdd

	err = service.setUserTotalPoint(userId, totalPoint)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) setUserRaundLevel(userId string, level int64) error {
	data := make(map[string]interface{})

	data[userId] = level

	_, err := service.redisClient.HMSet("user-raund-level", data)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getUserRaundLevel(userId string) (int, error) {
	var result string
	var err error
	result, err = service.redisClient.HMGet("user-raund-level", userId)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	level, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return level, nil
}

func (service *ColorService) getUserStageInfo(userId string, stage int) (int, error) {
	var result string
	var err error
	result, err = service.redisClient.HGet(fmt.Sprintf("user-stage-info-%s", userId), strconv.Itoa(stage))

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	level, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return level, nil
}

func (service *ColorService) setUserStageInfo(userId string, stage int, count int) error {

	_, err := service.redisClient.HSet(fmt.Sprintf("user-stage-info-%s", userId), strconv.Itoa(stage), count)

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) incrementUserStageInfo(userId string, stage int, count int) error {

	_, err := service.redisClient.HIncrBy(fmt.Sprintf("user-stage-info-%s", userId), strconv.Itoa(stage), int64(count))

	if err != nil {
		return types.NewBusinessException("system exception", "exp.systemexception")
	}

	return nil
}

func (service *ColorService) getDefaultStageInfo(stage int) (int, error) {
	var result string
	var err error
	result, err = service.redisClient.HMGet("default-stage-info", strconv.Itoa(stage))

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if result == "" {
		result = "0"
	}

	level, err := strconv.Atoi(result)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return level, nil
}
