package gin_release

import (
	"errors"
	"fmt"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
)

func ReleaseRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	engine := gin.New()

	engine.Use(gin.Logger(), CustomRecovery())
	return engine
}

func CustomRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {

		response.ErrorSystem(c, "", fmt.Sprintf("%s", err))
	})
}

type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	variable.ZapLog.Error(consts.ServerOccurredErrorMsg, zap.String("errStrace", errStr))
	return len(errStr), err
}
