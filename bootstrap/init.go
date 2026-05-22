package bootstrap

import (
	_ "ginskeleton/app/core/destroy"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/validator/common/register_validator"
	"ginskeleton/app/service/sys_log_hook"
	"ginskeleton/app/utils/casbin_v2"
	"ginskeleton/app/utils/gorm_v2"
	"ginskeleton/app/utils/snow_flake"
	"ginskeleton/app/utils/validator_translation"
	"ginskeleton/app/utils/websocket/core"
	"ginskeleton/app/utils/yml_config"
	"ginskeleton/app/utils/zap_factory"
	"log"
	"os"
)

func checkRequiredFolders() {
	if _, err := os.Stat(variable.BasePath + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
	if _, err := os.Stat(variable.BasePath + "/config/gorm_v2.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigGormNotExists + err.Error())
	}
	if _, err := os.Stat(variable.BasePath + "/public/"); err != nil {
		log.Fatal(my_errors.ErrorsPublicNotExists + err.Error())
	}
	if _, err := os.Stat(variable.BasePath + "/storage/logs/"); err != nil {
		if err := os.MkdirAll(variable.BasePath+"/storage/logs/", os.ModePerm); err != nil {
			log.Fatal(my_errors.ErrorsStorageLogsNotExists + err.Error())
		}
	}
	if _, err := os.Stat(variable.BasePath + "/public/storage"); os.IsNotExist(err) {
		os.MkdirAll(variable.BasePath+"/storage/app", os.ModePerm)
		os.Symlink(variable.BasePath+"/storage/app", variable.BasePath+"/public/storage")
	}
}

func init() {
	checkRequiredFolders()

	register_validator.WebRegisterValidator()
	register_validator.ApiRegisterValidator()

	variable.ConfigYml = yml_config.CreateYamlFactory()
	variable.ConfigYml.ConfigFileChangeListen()
	variable.ConfigGormv2Yml = variable.ConfigYml.Clone("gorm_v2")
	variable.ConfigGormv2Yml.ConfigFileChangeListen()

	variable.ZapLog = zap_factory.CreateZapFactory(sys_log_hook.ZapLogHandler)

	if variable.ConfigGormv2Yml.GetInt("Gormv2.Mysql.IsInitGlobalGormMysql") == 1 {
		if dbMysql, err := gorm_v2.GetOneMysqlClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.GormDbMysql = dbMysql
		}
	}
	if variable.ConfigGormv2Yml.GetInt("Gormv2.Sqlserver.IsInitGlobalGormSqlserver") == 1 {
		if dbSqlserver, err := gorm_v2.GetOneSqlserverClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.GormDbSqlserver = dbSqlserver
		}
	}
	if variable.ConfigGormv2Yml.GetInt("Gormv2.PostgreSql.IsInitGlobalGormPostgreSql") == 1 {
		if dbPostgre, err := gorm_v2.GetOnePostgreSqlClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.GormDbPostgreSql = dbPostgre
		}
	}

	variable.SnowFlake = snow_flake.CreateSnowflakeFactory()

	if variable.ConfigYml.GetInt("Websocket.Start") == 1 {
		variable.WebsocketHub = core.CreateHubFactory()
		if Wh, ok := variable.WebsocketHub.(*core.Hub); ok {
			go Wh.Run()
		}
	}

	if variable.ConfigYml.GetInt("Casbin.IsInit") == 1 {
		var err error
		if variable.Enforcer, err = casbin_v2.InitCasbinEnforcer(); err != nil {
			log.Fatal(err.Error())
		}
	}
	if err := validator_translation.InitTrans("zh"); err != nil {
		log.Fatal(my_errors.ErrorsValidatorTransInitFail + err.Error())
	}
}
