package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Login struct {

	BaseField
}

func (l Login) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&l); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(l, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "userLogin表单验证器json化失败", "")
	} else {

		(&web.Users{}).Login(extraAddBindDataContext)
	}

}
