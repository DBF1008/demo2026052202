package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Destroy struct {

	Id
}

func (d Destroy) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&d); err != nil {

		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(d, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserShow表单参数验证器json化失败", "")
		return
	} else {

		(&web.Users{}).Destroy(extraAddBindDataContext)

	}
}
