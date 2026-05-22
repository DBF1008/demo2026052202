package cur_userinfo

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/middleware/my_jwt"
	"github.com/gin-gonic/gin"
)

func GetCurrentUserId(context *gin.Context) (int64, bool) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	currentUser, exist := context.MustGet(tokenKey).(my_jwt.CustomClaims)
	return currentUser.UserId, exist
}
