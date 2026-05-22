package files

import (
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"mime/multipart"
	"net/http"
	"os"
)

func GetFilesMimeByFileName(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadOpenFail + err.Error())
	}
	defer f.Close()

	buffer := make([]byte, 32)
	if _, err := f.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	return http.DetectContentType(buffer)
}

func GetFilesMimeByFp(fp multipart.File) string {

	buffer := make([]byte, 32)
	if _, err := fp.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	return http.DetectContentType(buffer)
}
