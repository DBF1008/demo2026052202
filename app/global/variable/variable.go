package variable

import (
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/utils/snow_flake/snowflake_interf"
	"ginskeleton/app/utils/yml_config/ymlconfig_interf"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

var (
	BasePath           string
	EventDestroyPrefix = "Destroy_"
	ConfigKeyPrefix    = "Config_"
	DateFormat         = "2006-01-02 15:04:05"

	ZapLog *zap.Logger

	ConfigYml       ymlconfig_interf.YmlConfigInterf
	ConfigGormv2Yml ymlconfig_interf.YmlConfigInterf

	GormDbMysql      *gorm.DB
	GormDbSqlserver  *gorm.DB
	GormDbPostgreSql *gorm.DB

	SnowFlake snowflake_interf.InterfaceSnowFlake

	WebsocketHub              interface{}
	WebsocketHandshakeSuccess = `{"code":200,"msg":"ws连接成功","data":""}`
	WebsocketServerPingMsg    = "Server->Ping->Client"

	Enforcer *casbin.SyncedEnforcer

)

func init() {

	if curPath, err := os.Getwd(); err == nil {

		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
