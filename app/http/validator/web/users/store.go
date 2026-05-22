package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Store struct {
	BaseField

	RealName string `form:"real_name" json:"real_name" binding:"required,min=2"`
	Phone    string `form:"phone" json:"phone" binding:"required,len=11"`
	Remark   string `form:"remark" json:"remark" `
}

func (s Store) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&s); err != nil {

		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserStore表单验证器json化失败", "")
	} else {

		(&web.Users{}).Store(extraAddBindDataContext)
	}
}
