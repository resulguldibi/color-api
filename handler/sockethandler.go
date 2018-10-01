package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/service"
	"resulguldibi/color-api/types"
	"resulguldibi/color-api/util"

	"github.com/gin-gonic/gin"
)

func (handler SocketHandler) HandleAcceptMatchForMultiPlay(ctx *gin.Context, hub *service.SocketHub) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if !isExist {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	user = userData.(entity.User)

	handler.socketService.AcceptMatchForMultiPlay(ctx, user, hub)
}

func (handler SocketHandler) HandleRegisterForMultiPlay(ctx *gin.Context, hub *service.SocketHub) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if !isExist {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	user = userData.(entity.User)

	handler.socketService.RegisterForMultiPlay(ctx, user, hub)
}

func (handler SocketHandler) HandleUnRegisterForMultiPlay(ctx *gin.Context, hub *service.SocketHub) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if !isExist {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	user = userData.(entity.User)

	handler.socketService.UnRegisterForMultiPlay(ctx, user, hub)
}

func (handler SocketHandler) HandleMultiplayMove(ctx *gin.Context, hub *service.SocketHub) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if !isExist {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	request := &contract.MultiplayMatchMoveRequest{}

	if err := ctx.ShouldBindJSON(request); err == nil {

		user = userData.(entity.User)

		handler.socketService.MultiplayMove(ctx, request, user, hub)

	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}

}
