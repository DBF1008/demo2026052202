package web

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/model"
	"ginskeleton/app/service/users/curd"
	userstoken "ginskeleton/app/service/users/token"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
	"time"
)

type Users struct {
}

func (u *Users) Register(context *gin.Context) {

	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	userIp := context.ClientIP()
	if curd.CreateUserCurdFactory().Register(userName, pass, userIp) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdRegisterFailCode, consts.CurdRegisterFailMsg, "")
	}
}

func (u *Users) Login(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	phone := context.GetString(consts.ValidatorPrefix + "phone")
	userModelFact := model.CreateUserFactory("")
	userModel := userModelFact.Login(userName, pass)

	if userModel != nil {
		userTokenFactory := userstoken.CreateUserFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id, userModel.UserName, userModel.Phone, variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil {
			if userTokenFactory.RecordLoginToken(userToken, context.ClientIP()) {
				data := gin.H{
					"userId":     userModel.Id,
					"user_name":  userName,
					"realName":   userModel.RealName,
					"phone":      phone,
					"token":      userToken,
					"updated_at": time.Now().Format(variable.DateFormat),
				}
				response.Success(context, consts.CurdStatusOkMsg, data)
				go userModel.UpdateUserloginInfo(context.ClientIP(), userModel.Id)
				return
			}
		}
	}
	response.Fail(context, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, "")
}

func (u *Users) RefreshToken(context *gin.Context) {
	oldToken := context.GetString(consts.ValidatorPrefix + "token")
	if newToken, ok := userstoken.CreateUserFactory().RefreshToken(oldToken, context.ClientIP()); ok {
		res := gin.H{
			"token": newToken,
		}
		response.Success(context, consts.CurdStatusOkMsg, res)
	} else {
		response.Fail(context, consts.CurdRefreshTokenFailCode, consts.CurdRefreshTokenFailMsg, "")
	}
}

func (u *Users) Show(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit
	counts, showlist := model.CreateUserFactory("").Show(userName, int(limitStart), int(limit))
	if counts > 0 && showlist != nil {
		response.Success(context, consts.CurdStatusOkMsg, gin.H{"counts": counts, "list": showlist})
	} else {
		response.Fail(context, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

func (u *Users) Store(context *gin.Context) {
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	realName := context.GetString(consts.ValidatorPrefix + "real_name")
	phone := context.GetString(consts.ValidatorPrefix + "phone")
	remark := context.GetString(consts.ValidatorPrefix + "remark")

	if curd.CreateUserCurdFactory().Store(userName, pass, realName, phone, remark) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "")
	}
}

func (u *Users) Update(context *gin.Context) {

	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	userName := context.GetString(consts.ValidatorPrefix + "user_name")
	pass := context.GetString(consts.ValidatorPrefix + "pass")
	realName := context.GetString(consts.ValidatorPrefix + "real_name")
	phone := context.GetString(consts.ValidatorPrefix + "phone")
	remark := context.GetString(consts.ValidatorPrefix + "remark")
	userIp := context.ClientIP()

	if model.CreateUserFactory("").UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+", "+userName+" 已经被其他人使用", "")
		return
	}

	if curd.CreateUserCurdFactory().Update(int(userId), userName, pass, realName, phone, remark, userIp) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}

}

func (u *Users) Destroy(context *gin.Context) {

	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	if model.CreateUserFactory("").Destroy(int(userId)) {
		response.Success(context, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(context, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}
