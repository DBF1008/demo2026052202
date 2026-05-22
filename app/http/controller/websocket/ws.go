package websocket

import (
	serviceWs "ginskeleton/app/service/websocket"
	"github.com/gin-gonic/gin"
)

type Ws struct {
}

func (w *Ws) OnOpen(context *gin.Context) (*serviceWs.Ws, bool) {
	return (&serviceWs.Ws{}).OnOpen(context)
}

func (w *Ws) OnMessage(serviceWs *serviceWs.Ws, context *gin.Context) {
	serviceWs.OnMessage(context)
}
