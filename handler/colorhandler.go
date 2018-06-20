package handler

import (
	"fmt"
	"net/http"
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (handler ColorHandler) HandleGetRandomColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			responseSatus := util.PrepareResponseStatusWithMessage(false, fmt.Sprint(err))
			ctx.JSON(http.StatusBadRequest, responseSatus)
		}
	}()

	level, err := strconv.ParseInt(ctx.Request.URL.Query().Get("level"), 10, 64)
	if err != nil {
		panic(err)
	}

	response, err := handler.colorService.GetRandomColors(level)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}

func (handler ColorHandler) HandleValidateColors(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			responseSatus := util.PrepareResponseStatusWithMessage(false, fmt.Sprint(err))
			ctx.JSON(http.StatusBadRequest, responseSatus)
		}
	}()

	request := contract.ValidateColorsRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		response, err := handler.colorService.ValidateColors(request.SelectedColors, request.MixedColor)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		responseSatus := util.PrepareResponseStatusWithMessage(false, err.Error())
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}
