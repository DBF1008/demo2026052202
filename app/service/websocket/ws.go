package websocket

import (
	"fmt"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/websocket/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Ws struct {
	WsClient *core.Client
}

func (w *Ws) OnOpen(context *gin.Context) (*Ws, bool) {
	if client, ok := (&core.Client{}).OnOpen(context); ok {

		token := context.GetString(consts.ValidatorPrefix + "token")
		variable.ZapLog.Info("获取到的客户端上线时携带的唯一标记值：", zap.String("token", token))

		w.WsClient = client
		go w.WsClient.Heartbeat()
		return w, true
	} else {
		return nil, false
	}
}

func (w *Ws) OnMessage(context *gin.Context) {
	go w.WsClient.ReadPump(func(messageType int, receivedData []byte) {

		tempMsg := "服务器已经收到了你的消息==>" + string(receivedData)

		if err := w.WsClient.SendMessage(messageType, tempMsg); err != nil {
			variable.ZapLog.Error("消息发送出现错误", zap.Error(err))
		}

	}, w.OnError, w.OnClose)
}

func (w *Ws) OnError(err error) {
	w.WsClient.State = 0
	variable.ZapLog.Error("远端掉线、卡死、刷新浏览器等会触发该错误:", zap.Error(err))

}

func (w *Ws) OnClose() {

	w.WsClient.Hub.UnRegister <- w.WsClient
}

func (w *Ws) GetOnlineClients() {

	fmt.Printf("在线客户端数量：%d\n", len(w.WsClient.Hub.Clients))
}

func (w *Ws) BroadcastMsg(sendMsg string) {
	for onlineClient := range w.WsClient.Hub.Clients {

		if err := onlineClient.SendMessage(websocket.TextMessage, sendMsg); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
		}
	}
}
