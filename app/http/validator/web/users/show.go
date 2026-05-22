package users

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/web"
	common_data_type "ginskeleton/app/http/validator/common/data_type"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Show struct {

	UserName string `form:"user_name" json:"user_name"  binding:"required,min=1"`
	common_data_type.Page
}

func (s Show) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&s); err != nil {

		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(s, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "UserShow表单验证器json化失败", "")
	} else {

		(&web.Users{}).Show(extraAddBindDataContext)
	}
}
