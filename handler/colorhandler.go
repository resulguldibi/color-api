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

func (handler ColorHandler) HandleGetRandomColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			exp := &types.ExceptionMessage{}
			_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
			responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
			ctx.JSON(http.StatusBadRequest, responseSatus)
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

func (handler ColorHandler) HandleValidateColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			exp := &types.ExceptionMessage{}
			_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
			responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
			ctx.JSON(http.StatusBadRequest, responseSatus)
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
		response, err := handler.colorService.ValidateColors(user.Id, key, request.SelectedColors, request.MixedColor)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}
