package upload_file

import (
	"errors"
	"fmt"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/md5_encrypt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
	"time"
)

func Upload(context *gin.Context, savePath string) (r bool, finnalSavePath interface{}) {

	newSavePath, newReturnPath := generateYearMonthPath(savePath)

	file, _ := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField"))

	var saveErr error
	if sequence := variable.SnowFlake.GetId(); sequence > 0 {
		saveFileName := fmt.Sprintf("%d%s", sequence, file.Filename)
		saveFileName = md5_encrypt.MD5(saveFileName) + path.Ext(saveFileName)

		if saveErr = context.SaveUploadedFile(file, newSavePath+saveFileName); saveErr == nil {

			finnalSavePath = gin.H{
				"path": strings.ReplaceAll(newReturnPath+saveFileName, variable.BasePath, ""),
			}
			return true, finnalSavePath
		}
	} else {
		saveErr = errors.New(my_errors.ErrorsSnowflakeGetIdFail)
		variable.ZapLog.Error("文件保存出错：" + saveErr.Error())
	}
	return false, nil

}

func generateYearMonthPath(savePathPre string) (string, string) {
	returnPath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileReturnPath")
	curYearMonth := time.Now().Format("2006_01")
	newSavePathPre := savePathPre + curYearMonth
	newReturnPathPre := returnPath + curYearMonth

	if _, err := os.Stat(newSavePathPre); err != nil {
		if err = os.MkdirAll(newSavePathPre, os.ModePerm); err != nil {
			variable.ZapLog.Error("文件上传创建目录出错" + err.Error())
			return "", ""
		}
	}
	return newSavePathPre + "/", newReturnPathPre + "/"
}
