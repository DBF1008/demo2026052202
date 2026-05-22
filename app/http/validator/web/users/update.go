package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Update struct {
	BaseField
	Id

	RealName string `form:"real_name" json:"real_name" binding:"required,min=2"`
	Phone    string `form:"phone" json:"phone" binding:"required,len=11"`
	Remark   string `form:"remark" json:"remark"`
}

func (u Update) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&u); err != nil {

		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(u, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserUpdate表单验证器json化失败", "")
	} else {

		(&web.Users{}).Update(extraAddBindDataContext)
	}
}
