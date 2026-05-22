package core

import (
	"errors"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/service/websocket/on_open_success"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Hub                *Hub
	Conn               *websocket.Conn
	Send               chan []byte
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailTimes int
	ClientLastPongTime time.Time
	State              uint8
	sync.RWMutex
	on_open_success.ClientMoreParams
}

func (c *Client) OnOpen(context *gin.Context) (*Client, bool) {

	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				variable.ZapLog.Error(my_errors.ErrorsWebsocketOnOpenFail, zap.Error(val))
			}
		}
	}()
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  variable.ConfigYml.GetInt("Websocket.WriteReadBufferSize"),
		WriteBufferSize: variable.ConfigYml.GetInt("Websocket.WriteReadBufferSize"),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if wsConn, err := upGrader.Upgrade(context.Writer, context.Request, nil); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsWebsocketUpgradeFail + err.Error())
		return nil, false
	} else {
		if wsHub, ok := variable.WebsocketHub.(*Hub); ok {
			c.Hub = wsHub
		}
		c.Conn = wsConn
		c.Send = make(chan []byte, variable.ConfigYml.GetInt("Websocket.WriteReadBufferSize"))
		c.PingPeriod = time.Second * variable.ConfigYml.GetDuration("Websocket.PingPeriod")
		c.ReadDeadline = time.Second * variable.ConfigYml.GetDuration("Websocket.ReadDeadline")
		c.WriteDeadline = time.Second * variable.ConfigYml.GetDuration("Websocket.WriteDeadline")

		if err := c.SendMessage(websocket.TextMessage, variable.WebsocketHandshakeSuccess); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsWebsocketWriteMgsFail, zap.Error(err))
		}
		c.Conn.SetReadLimit(variable.ConfigYml.GetInt64("Websocket.MaxMessageSize"))
		c.Hub.Register <- c
		c.State = 1
		c.ClientLastPongTime = time.Now()
		return c, true
	}

}

func (c *Client) ReadPump(callbackOnMessage func(messageType int, receivedData []byte), callbackOnError func(err error), callbackOnClose func()) {

	defer func() {
		err := recover()
		if err != nil {
			if realErr, isOk := err.(error); isOk {
				variable.ZapLog.Error(my_errors.ErrorsWebsocketReadMessageFail, zap.Error(realErr))
			}
		}
		callbackOnClose()
	}()

	for {
		if c.State == 1 {
			mt, bReceivedData, err := c.Conn.ReadMessage()
			if err == nil {
				callbackOnMessage(mt, bReceivedData)
			} else {

				callbackOnError(err)
				break
			}
		} else {

			callbackOnError(errors.New(my_errors.ErrorsWebsocketStateInvalid))
			break
		}

	}
}

func (c *Client) SendMessage(messageType int, message string) error {
	c.Lock()
	defer func() {
		c.Unlock()
	}()

	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsWebsocketSetWriteDeadlineFail, zap.Error(err))
		return err
	}
	if err := c.Conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Client) Heartbeat() {

	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				variable.ZapLog.Error(my_errors.ErrorsWebsocketBeatHeartFail, zap.Error(val))
			}
		}
		ticker.Stop()
	}()

	if c.ReadDeadline == 0 {
		_ = c.Conn.SetReadDeadline(time.Time{})
	} else {
		_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
	}
	c.Conn.SetPongHandler(func(receivedPong string) error {
		if c.ReadDeadline > time.Nanosecond {
			_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		} else {
			_ = c.Conn.SetReadDeadline(time.Time{})
		}

		c.ClientLastPongTime = time.Now()

		return nil
	})

	for {
		select {
		case <-ticker.C:
			if c.State == 1 {

				serverAllowMaxOfflineSeconds := float64(variable.ConfigYml.GetInt("Websocket.HeartbeatFailMaxTimes")) * (float64(variable.ConfigYml.GetDuration("Websocket.PingPeriod")))
				if time.Now().Sub(c.ClientLastPongTime).Seconds() > serverAllowMaxOfflineSeconds {
					c.State = 0
					c.Hub.UnRegister <- c
					variable.ZapLog.Warn(my_errors.ErrorsWebsocketClientOfflineTimeout, zap.Float64("timeout(seconds): ", serverAllowMaxOfflineSeconds))
					return
				}

				if err := c.SendMessage(websocket.PingMessage, variable.WebsocketServerPingMsg); err != nil {
					c.HeartbeatFailTimes++
					if c.HeartbeatFailTimes > variable.ConfigYml.GetInt("Websocket.HeartbeatFailMaxTimes") {
						c.State = 0
						c.Hub.UnRegister <- c
						variable.ZapLog.Error(my_errors.ErrorsWebsocketBeatHeartsMoreThanMaxTimes, zap.Error(err))
						return
					}
				} else {
					if c.HeartbeatFailTimes > 0 {
						c.HeartbeatFailTimes--
					}
				}
			} else {
				return
			}

		}
	}
}
