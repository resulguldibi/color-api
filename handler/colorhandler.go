package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/types"
	"resulguldibi/color-api/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler ColorHandler) HandleRaundHistory(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if isExist {
		user = userData.(entity.User)
	} else {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	key := ctx.GetHeader("RaundKey")

	if key == "" {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	response, err := handler.colorService.GetUserRaundHistory(user.Id, key)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleColorStepHelp(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if isExist {
		user = userData.(entity.User)
	} else {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	key := ctx.GetHeader("RaundKey")

	if key == "" {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	colorsJson := ctx.Request.URL.Query().Get("colors")
	colors := []*entity.Color{}
	err := json.Unmarshal([]byte(colorsJson), &colors)

	if err != nil {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	response, err := handler.colorService.GetColorStepHelp(user.Id, key, colors)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleColorHelp(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if isExist {
		user = userData.(entity.User)
	} else {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	key := ctx.GetHeader("RaundKey")

	if key == "" {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	response, err := handler.colorService.GetColorHelp(user.Id, key)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleGetRandomColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	level, err := strconv.ParseInt(ctx.Request.URL.Query().Get("level"), 10, 64)
	if err != nil {
		panic(err)
	}

	userData, isExist := ctx.Get("User")
	var user entity.User
	if isExist {
		user = userData.(entity.User)
	}

	response, err := handler.colorService.GetRandomColors(user.Id, level)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleGetLevels(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	response, err := handler.colorService.GetLevels()
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleValidateColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.ValidateColorsRequest{}
	key := ctx.GetHeader("RaundKey")

	if err := ctx.ShouldBindJSON(&request); err == nil {
		userData, isExist := ctx.Get("User")
		var user entity.User
		if isExist {
			user = userData.(entity.User)
		}
		response, err := handler.colorService.ValidateColors(user.Id, key, request.SelectedColors)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler ColorHandler) HandleRankings(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	key := ctx.GetHeader("RaundKey")

	userData, isExist := ctx.Get("User")
	var user entity.User
	if isExist {
		user = userData.(entity.User)
	}
	response, err := handler.colorService.GetRankings(user.Id, key)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}
