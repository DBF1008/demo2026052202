package upload_files

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/controller/web"
	"ginskeleton/app/utils/files"
	"ginskeleton/app/utils/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type UpFiles struct {
}

func (u UpFiles) CheckParams(context *gin.Context) {
	tmpFile, err := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField"))
	var isPass bool

	if err != nil {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, err.Error())
		return
	}
	if tmpFile.Size == 0 {
		response.Fail(context, consts.FilesUploadMoreThanMaxSizeCode, consts.FilesUploadIsEmpty, "")
		return
	}

	sizeLimit := variable.ConfigYml.GetInt64("FileUploadSetting.Size")
	if tmpFile.Size > sizeLimit<<20 {
		response.Fail(context, consts.FilesUploadMoreThanMaxSizeCode, consts.FilesUploadMoreThanMaxSizeMsg+strconv.FormatInt(sizeLimit, 10)+"M", "")
		return
	}

	if fp, err := tmpFile.Open(); err == nil {
		mimeType := files.GetFilesMimeByFp(fp)

		for _, value := range variable.ConfigYml.GetStringSlice("FileUploadSetting.AllowMimeType") {
			if strings.ReplaceAll(value, " ", "") == strings.ReplaceAll(mimeType, " ", "") {
				isPass = true
				break
			}
		}
		_ = fp.Close()
	} else {
		response.ErrorSystem(context, consts.ServerOccurredErrorMsg, "")
		return
	}

	if !isPass {
		response.Fail(context, consts.FilesUploadMimeTypeFailCode, consts.FilesUploadMimeTypeFailMsg, "")
	} else {
		(&web.Upload{}).StartUpload(context)
	}
}
