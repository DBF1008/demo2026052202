package register_validator

import (
	"ginskeleton/app/core/container"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/validator/common/upload_files"
	"ginskeleton/app/http/validator/common/websocket"
	"ginskeleton/app/http/validator/web/users"
)

func WebRegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = consts.ValidatorPrefix + "UsersRegister"
	containers.Set(key, users.Register{})
	key = consts.ValidatorPrefix + "UsersLogin"
	containers.Set(key, users.Login{})
	key = consts.ValidatorPrefix + "RefreshToken"
	containers.Set(key, users.RefreshToken{})

	key = consts.ValidatorPrefix + "UsersShow"
	containers.Set(key, users.Show{})
	key = consts.ValidatorPrefix + "UsersStore"
	containers.Set(key, users.Store{})
	key = consts.ValidatorPrefix + "UsersUpdate"
	containers.Set(key, users.Update{})
	key = consts.ValidatorPrefix + "UsersDestroy"
	containers.Set(key, users.Destroy{})

	key = consts.ValidatorPrefix + "UploadFiles"
	containers.Set(key, upload_files.UpFiles{})

	key = consts.ValidatorPrefix + "WebsocketConnect"
	containers.Set(key, websocket.Connect{})
}
