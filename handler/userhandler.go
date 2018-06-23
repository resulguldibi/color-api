package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resulguldibi/color-api/types"
	"resulguldibi/color-api/util"

	"github.com/gin-gonic/gin"
)

func (handler UserHandler) HandleSignIn(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			exp := &types.ExceptionMessage{}
			_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
			responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
			ctx.JSON(http.StatusBadRequest, responseSatus)
		}
	}()

	response, err := handler.userService.SignIn(1)
	util.CheckErr(err)
	ctx.JSON(http.StatusOK, response)
}
