package home

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/controller/api"
	common_data_type "ginskeleton/app/http/validator/common/data_type"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type News struct {
	NewsType string `form:"newsType" json:"newsType"  binding:"required,min=1"`
	common_data_type.Page
}

func (n News) CheckParams(context *gin.Context) {

	if err := context.ShouldBind(&n); err != nil {

		response.ValidatorError(context, err)
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "HomeNews表单验证器json化失败", "")
	} else {

		(&api.Home{}).News(extraAddBindDataContext)
	}

}
