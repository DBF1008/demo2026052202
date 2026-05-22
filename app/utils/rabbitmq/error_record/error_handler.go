package error_record

import "ginskeleton/app/global/variable"

func ErrorDeal(err error) error {
	if err != nil {
		variable.ZapLog.Error(err.Error())
	}
	return err
}
