package authorization

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	userstoken "ginskeleton/app/service/users/token"
	"ginskeleton/app/utils/response"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

func CheckTokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}

		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.TokenErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			tokenIsEffective := userstoken.CreateUserFactory().IsEffective(token[1])
			if tokenIsEffective {
				if customToken, err := userstoken.CreateUserFactory().ParseToken(token[1]); err == nil {
					key := variable.ConfigYml.GetString("Token.BindContextKeyName")

					context.Set(key, customToken)
				}
				context.Next()
			} else {
				response.ErrorTokenAuthFail(context)
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

func CheckTokenAuthWithRefresh() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}

		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.TokenErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			tokenIsEffective := userstoken.CreateUserFactory().IsEffective(token[1])

			if tokenIsEffective {
				if customToken, err := userstoken.CreateUserFactory().ParseToken(token[1]); err == nil {
					key := variable.ConfigYml.GetString("Token.BindContextKeyName")

					context.Set(key, customToken)

					context.Header("Refresh-Token", "")
					context.Header("Access-Control-Expose-Headers", "Refresh-Token")
				}
				context.Next()
			} else {

				if userstoken.CreateUserFactory().TokenIsMeetRefreshCondition(token[1]) {

					if newToken, ok := userstoken.CreateUserFactory().RefreshToken(token[1], context.ClientIP()); ok {
						if customToken, err := userstoken.CreateUserFactory().ParseToken(newToken); err == nil {
							key := variable.ConfigYml.GetString("Token.BindContextKeyName")

							context.Set(key, customToken)
						}

						context.Header("Refresh-Token", newToken)
						context.Header("Access-Control-Expose-Headers", "Refresh-Token")
						context.Next()
					} else {
						response.ErrorTokenRefreshFail(context)
					}
				} else {
					response.ErrorTokenRefreshFail(context)
				}
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

func RefreshTokenConditionCheck() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.TokenErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {

			if userstoken.CreateUserFactory().TokenIsMeetRefreshCondition(token[1]) {
				context.Next()
			} else {
				response.ErrorTokenRefreshFail(context)
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

func CheckCasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requstUrl := c.Request.URL.Path
		method := c.Request.Method

		role := "2"

		isPass, err := variable.Enforcer.Enforce(role, requstUrl, method)
		if err != nil {
			response.ErrorCasbinAuthFail(c, err.Error())
			return
		} else if !isPass {
			response.ErrorCasbinAuthFail(c, "")
			return
		} else {
			c.Next()
		}
	}
}

func CheckCaptchaAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
		captchaValueKey := variable.ConfigYml.GetString("Captcha.captchaValue")
		captchaId := c.PostForm(captchaIdKey)
		value := c.PostForm(captchaValueKey)
		if captchaId == "" || value == "" {
			response.Fail(c, consts.CaptchaCheckParamsInvalidCode, consts.CaptchaCheckParamsInvalidMsg, "")
			return
		}
		if captcha.VerifyString(captchaId, value) {
			c.Next()
		} else {
			response.Fail(c, consts.CaptchaCheckFailCode, consts.CaptchaCheckFailMsg, "")
		}
	}
}
