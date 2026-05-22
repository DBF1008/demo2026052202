package websocket

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	controllerWs "ginskeleton/app/http/controller/websocket"
	"ginskeleton/app/http/validator/core/data_transfer"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Connect struct {
	Token string `form:"token" json:"token" binding:"required,min=10"`
}

func (c Connect) CheckParams(context *gin.Context) {

	if variable.ConfigYml.GetInt("Websocket.Start") != 1 {
		variable.ZapLog.Error(consts.WsServerNotStartMsg)
		return
	}

	if err := context.ShouldBind(&c); err != nil {
		variable.ZapLog.Error("客户端上线参数不合格", zap.Error(err))
		return
	}
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		variable.ZapLog.Error("websocket-Connect 表单验证器json化失败")
		context.Abort()
		return
	} else {
		if serviceWs, ok := (&controllerWs.Ws{}).OnOpen(extraAddBindDataContext); ok == false {
			variable.ZapLog.Error(consts.WsOpenFailMsg)
		} else {
			(&controllerWs.Ws{}).OnMessage(serviceWs, extraAddBindDataContext)
		}
	}
}
