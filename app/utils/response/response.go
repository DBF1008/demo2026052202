package response

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/utils/validator_translation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func ReturnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

func ReturnJsonFromString(Context *gin.Context, httpCode int, jsonStr string) {
	Context.Header("Content-Type", "application/json; charset=utf-8")
	Context.String(httpCode, jsonStr)
}

func Success(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, msg, data)
}

func Fail(c *gin.Context, dataCode int, msg string, data interface{}) {
	ReturnJson(c, http.StatusBadRequest, dataCode, msg, data)
	c.Abort()
}

func ErrorTokenBaseInfo(c *gin.Context) {
	ReturnJson(c, http.StatusBadRequest, http.StatusBadRequest, my_errors.ErrorsTokenBaseInfo, "")

	c.Abort()
}

func ErrorTokenAuthFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsNoAuthorization, "")

	c.Abort()
}

func ErrorTokenRefreshFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsRefreshTokenFail, "")

	c.Abort()
}

func TokenErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusUnauthorized, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

func ErrorCasbinAuthFail(c *gin.Context, msg interface{}) {
	ReturnJson(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, my_errors.ErrorsCasbinNoAuthorization, msg)
	c.Abort()
}

func ErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.ServerOccurredErrorCode, consts.ServerOccurredErrorMsg+msg, data)
	c.Abort()
}

func ValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		wrongParam := validator_translation.RemoveTopStruct(errs.Translate(validator_translation.Trans))
		ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	} else {
		errStr := err.Error()

		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": my_errors.ErrorNotAllParamsIsBlank})
		} else {
			ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, gin.H{"tips": errStr})
		}
	}
	c.Abort()
}
