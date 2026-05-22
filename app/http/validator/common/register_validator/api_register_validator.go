package register_validator

import (
	"ginskeleton/app/core/container"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/validator/api/home"
)

func ApiRegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = consts.ValidatorPrefix + "HomeNews"
	containers.Set(key, home.News{})
}
