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

func (service *ColorService) GetRandomColors(userId string, level int64) (*contract.GetColorResponse, error) {
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

	code := util.GenerateGuid()
	response.MixedColor = mixedColor
	response.RandomColors = finalRandomColors
	response.Code = code

	//save generated random colors to use in /validate endpoint
	allColors := append(finalRandomColors, mixedColor)
	err = service.setUserRaundGeneratedRandomColors(code, allColors)
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
		err = service.deleteUserExistingRaundGeneratedRandomColors(existingKey)
		if err != nil {
			panic(err)
		}
	}

	err = service.setUserRaundKey(userId, code)
	if err != nil {
		panic(err)
	}

	err = service.setUserRaundLevel(userId, level)

	if err != nil {
		panic(err)
	}

	return response, err
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

func (service *ColorService) getUserRaundKey(userId string) (string, error) {

	response, err := service.redisClient.HMGet("user-raund-key", userId)

	if err != nil {
		return response, types.NewBusinessException("system exception", "exp.systemexception")
	}

	return response, err
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

func (service *ColorService) deleteUserExistingRaundGeneratedRandomColors(key string) error {

	_, err := service.redisClient.HDel("user-raund-generated-random-colors", key)

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

func (service *ColorService) calculateUserRaundPoint(userId string, key string) (int, error) {

	//hmset user-raund-point "1234" "12s12-12sas-3asw12-12sa1" "12s12-12sas-3asw12-12sa1" "100"
	step, err := service.getUserRaundStepNumber(userId, key)

	//raund full point is level * 20

	var level int
	level, err = service.getUserRaundLevel(userId)

	if err != nil {
		return 0, types.NewBusinessException("system exception", "exp.systemexception")
	}

	raundPoint := level * raundStartPoint

	raundPoint = raundPoint - step

	return raundPoint, nil
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

func (service *ColorService) ValidateColors(userId string, sendedKey string, colors []*entity.Color, color *entity.Color) (*contract.ValidateColorsResponse, error) {
	response := &contract.ValidateColorsResponse{}
	var err error
	var key string

	var isColorsValid bool = false

	// get code from client (code was sended to client in /colors response) in /validate request to calculate user point (client should send this guid in /validate request)
	// validate raund key

	key, err = service.getUserRaundKey(userId)

	if err != nil {
		return nil, types.NewBusinessException("system exception", "exp.systemexception")
	}

	if key != sendedKey {
		return nil, types.NewBusinessException("invalid key", "exp.invalidkey")
	}

	// validate sended color
	allColors := append(colors, color)
	isColorsValid, err = service.validateSendedColors(userId, key, allColors)

	if err != nil {
		return nil, err
	}

	if !isColorsValid {
		return nil, types.NewBusinessException("invalid colors", "exp.invalidcolors")
	}

	// if step number is reached to max retry number, then game is over.

	err = service.checkUserRaundStepNumber(userId, key)
	if err != nil {
		return nil, err
	}
	// increment user step number in every /validate request (store this data in redis        -> hmset user-step-number "1234" "12s12-12sas-3asw12-12sa1" "12s12-12sas-3asw12-12sa1" "10")

	err = service.updateUserRaundStepNumber(userId, key)
	if err != nil {
		return nil, err
	}

	isMatched := util.IsMatchedColors(colors, color)

	if !isMatched {
		// if step number is reached to max retry number, then game is over.

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
	} else {
		// calculate point with generated point algorithm.
		raundPoint, err := service.calculateUserRaundPoint(userId, key)
		if err != nil {
			return nil, err
		}
		// update user point
		raundPoint = raundPoint + 1
		err = service.setUserRaundPoint(userId, key, raundPoint)
		if err != nil {
			return nil, err
		}

		err = service.updateUserTotalPoint(userId, raundPoint)
		if err != nil {
			return nil, err
		}
	}

	response.IsValid = isMatched
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
