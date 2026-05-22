package yml_config

import (
	"ginskeleton/app/core/container"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/yml_config/ymlconfig_interf"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

var lastChangeTime time.Time
var containerFactory = container.CreateContainersFactory()

func init() {
	lastChangeTime = time.Now()
}

func CreateYamlFactory(fileName ...string) ymlconfig_interf.YmlConfigInterf {

	yamlConfig := viper.New()

	yamlConfig.AddConfigPath(variable.BasePath + "/config")

	if len(fileName) == 0 {
		yamlConfig.SetConfigName("config")
	} else {
		yamlConfig.SetConfigName(fileName[0])
	}

	yamlConfig.SetConfigType("yml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(my_errors.ErrorsConfigInitFail + err.Error())
	}

	return &ymlConfig{
		viper: yamlConfig,
		mu:    new(sync.Mutex),
	}
}

type ymlConfig struct {
	viper *viper.Viper
	mu    *sync.Mutex
}

func (y *ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

func (y *ymlConfig) keyIsCache(keyName string) bool {
	if _, exists := containerFactory.KeyIsExists(variable.ConfigKeyPrefix + keyName); exists {
		return true
	} else {
		return false
	}
}

func (y *ymlConfig) cache(keyName string, value interface{}) bool {

	y.mu.Lock()
	defer y.mu.Unlock()
	if _, exists := containerFactory.KeyIsExists(variable.ConfigKeyPrefix + keyName); exists {
		return true
	}
	return containerFactory.Set(variable.ConfigKeyPrefix+keyName, value)
}

func (y *ymlConfig) getValueFromCache(keyName string) interface{} {
	return containerFactory.Get(variable.ConfigKeyPrefix + keyName)
}

func (y *ymlConfig) clearCache() {
	containerFactory.FuzzyDelete(variable.ConfigKeyPrefix)
}

func (y *ymlConfig) Clone(fileName string) ymlconfig_interf.YmlConfigInterf {

	var ymlC = *y
	var ymlConfViper = *(y.viper)
	(&ymlC).viper = &ymlConfViper

	(&ymlC).viper.SetConfigName(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {

		if variable.ZapLog == nil {
			log.Fatal(my_errors.ErrorsConfigInitFail+", error: ", err.Error())
		} else {
			variable.ZapLog.Error(my_errors.ErrorsConfigInitFail, zap.Error(err))
		}
	}
	return &ymlC
}

func (y *ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}

}

func (y *ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y *ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName) {
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}
