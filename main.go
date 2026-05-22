package main

import (
	"ginskeleton/app/global/variable"
	_ "ginskeleton/bootstrap"
	"ginskeleton/routers"
)

func main() {
	router := routers.InitWebRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Web.Port"))
}