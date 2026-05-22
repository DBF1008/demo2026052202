package api

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Home struct {
}

func (u *Home) News(context *gin.Context) {

	newsType := context.GetString(consts.ValidatorPrefix + "newsType")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	userIp := context.ClientIP()
	ref := context.GetHeader("Referer")

	response.Success(context, "ok", gin.H{
		"newsType": newsType,
		"page":     page,
		"limit":    limit,
		"userIp":   userIp,
		"title":    "门户首页公司新闻标题001",
		"content":  "门户新闻内容001",
		"referer":  ref,
	})
}
