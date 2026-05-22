package routers

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/controller/captcha"
	"ginskeleton/app/http/middleware/authorization"
	"ginskeleton/app/http/middleware/cors"
	validatorFactory "ginskeleton/app/http/validator/core/factory"
	"ginskeleton/app/utils/gin_release"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitWebRouter() *gin.Engine {
	var router *gin.Engine

	if variable.ConfigYml.GetBool("AppDebug") == false {

		router = gin_release.ReleaseRouter()

	} else {

		router = gin.Default()
		pprof.Register(router)
	}

	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "HelloWorld,这是后端模块")
	})

	router.Static("/public", "./public")
	router.StaticFS("/dir", http.Dir("./public"))
	router.StaticFile("/abcd", "./public/readme.md")

	verifyCode := router.Group("captcha")
	{

		verifyCode.GET("/", (&captcha.Captcha{}).GenerateId)
		verifyCode.GET("/:captcha_id", (&captcha.Captcha{}).GetImg)
		verifyCode.GET("/:captcha_id/:captcha_value", (&captcha.Captcha{}).CheckCode)
	}

	backend := router.Group("/admin/")
	{

		backend.GET("ws", validatorFactory.Create(consts.ValidatorPrefix+"WebsocketConnect"))

		noAuth := backend.Group("users/")
		{

			noAuth.POST("register", validatorFactory.Create(consts.ValidatorPrefix+"UsersRegister"))

			noAuth.POST("login", validatorFactory.Create(consts.ValidatorPrefix+"UsersLogin"))

		}

		refreshToken := backend.Group("users/")
		{

			refreshToken.Use(authorization.RefreshTokenConditionCheck()).POST("refreshtoken", validatorFactory.Create(consts.ValidatorPrefix+"RefreshToken"))
		}

		backend.Use(authorization.CheckTokenAuth())
		{

			users := backend.Group("users/")
			{

				users.GET("index", validatorFactory.Create(consts.ValidatorPrefix+"UsersShow"))

				users.POST("create", validatorFactory.Create(consts.ValidatorPrefix+"UsersStore"))

				users.POST("edit", validatorFactory.Create(consts.ValidatorPrefix+"UsersUpdate"))

				users.POST("delete", validatorFactory.Create(consts.ValidatorPrefix+"UsersDestroy"))
			}

			uploadFiles := backend.Group("upload/")
			{
				uploadFiles.POST("files", validatorFactory.Create(consts.ValidatorPrefix+"UploadFiles"))
			}
		}
	}
	return router
}
