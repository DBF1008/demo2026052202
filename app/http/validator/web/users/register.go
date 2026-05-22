package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Register struct {
	BaseField

	Phone  string `form:"phone" json:"phone"`
	CardNo string `form:"card_no" json:"card_no"`
}

func (r Register) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&r); err != nil {
		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(r, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserRegister表单验证器json化失败", "")
	} else {

		(&web.Users{}).Register(extraAddBindDataContext)
	}

}
