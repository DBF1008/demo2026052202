package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"github.com/gin-gonic/gin"
)

type DestroyAfter struct{}

func (d *DestroyAfter) After(context *gin.Context) {

	go func() {
		userId := context.GetFloat64(consts.ValidatorPrefix + "id")
		variable.ZapLog.Sugar().Infof("模拟 Users 删除操作， After 回调,用户ID：%.f\n", userId)
	}()
}
