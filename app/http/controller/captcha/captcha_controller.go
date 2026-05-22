package captcha

import (
	"bytes"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/response"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
)

type Captcha struct{}

func (c *Captcha) GenerateId(context *gin.Context) {

	var length = variable.ConfigYml.GetInt("Captcha.length")
	var captchaId, imgUrl, refresh, verify string

	captchaId = captcha.NewLen(length)
	imgUrl = "/captcha/" + captchaId + ".png"
	refresh = imgUrl + "?reload=1"
	verify = "/captcha/" + captchaId + "/这里替换为正确的验证码进行验证"

	response.Success(context, "验证码信息", gin.H{
		"id":      captchaId,
		"img_url": imgUrl,
		"refresh": refresh,
		"verify":  verify,
	})

}

func (c *Captcha) GetImg(context *gin.Context) {
	captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
	captchaId := context.Param(captchaIdKey)
	_, file := path.Split(context.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || captchaId == "" {
		response.Fail(context, consts.CaptchaGetParamsInvalidCode, consts.CaptchaGetParamsInvalidMsg, "")
		return
	}

	if context.Query("reload") != "" {
		captcha.Reload(id)
	}

	context.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Header("Pragma", "no-cache")
	context.Header("Expires", "0")

	var vBytes bytes.Buffer
	if ext == ".png" {
		context.Header("Content-Type", "image/png")

		_ = captcha.WriteImage(&vBytes, id, captcha.StdWidth, captcha.StdHeight)
		http.ServeContent(context.Writer, context.Request, id+ext, time.Time{}, bytes.NewReader(vBytes.Bytes()))
	}
}

func (c *Captcha) CheckCode(context *gin.Context) {
	captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
	captchaValueKey := variable.ConfigYml.GetString("Captcha.captchaValue")

	captchaId := context.Param(captchaIdKey)
	value := context.Param(captchaValueKey)

	if captchaId == "" || value == "" {
		response.Fail(context, consts.CaptchaCheckParamsInvalidCode, consts.CaptchaCheckParamsInvalidMsg, "")
		return
	}
	if captcha.VerifyString(captchaId, value) {
		response.Success(context, consts.CaptchaCheckOkMsg, "")
	} else {
		response.Fail(context, consts.CaptchaCheckFailCode, consts.CaptchaCheckFailMsg, "")
	}
}
