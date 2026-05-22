package gorm_v2

import (
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/global/variable"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func MaskNotDataError(gormDB *gorm.DB) {
	gormDB.Statement.RaiseErrorOnNotFound = false
}

func CreateBeforeHook(gormDB *gorm.DB) {
	if reflect.TypeOf(gormDB.Statement.Dest).Kind() != reflect.Ptr {
		variable.ZapLog.Warn(my_errors.ErrorsGormDBCreateParamsNotPtr)
	} else {
		destValueOf := reflect.ValueOf(gormDB.Statement.Dest).Elem()
		if destValueOf.Type().Kind() == reflect.Slice || destValueOf.Type().Kind() == reflect.Array {
			inLen := destValueOf.Len()
			for i := 0; i < inLen; i++ {
				row := destValueOf.Index(i)
				if row.Type().Kind() == reflect.Struct {
					if b, column := structHasSpecialField("CreatedAt", row); b {
						destValueOf.Index(i).FieldByName(column).Set(reflect.ValueOf(time.Now().Format(variable.DateFormat)))
					}
					if b, column := structHasSpecialField("UpdatedAt", row); b {
						destValueOf.Index(i).FieldByName(column).Set(reflect.ValueOf(time.Now().Format(variable.DateFormat)))
					}

				} else if row.Type().Kind() == reflect.Map {
					if b, column := structHasSpecialField("created_at", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(variable.DateFormat)))
					}
					if b, column := structHasSpecialField("updated_at", row); b {
						row.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(variable.DateFormat)))
					}
				}
			}
		} else if destValueOf.Type().Kind() == reflect.Struct {

			if b, column := structHasSpecialField("CreatedAt", gormDB.Statement.Dest); b {
				gormDB.Statement.SetColumn(column, time.Now().Format(variable.DateFormat))
			}
			if b, column := structHasSpecialField("UpdatedAt", gormDB.Statement.Dest); b {
				gormDB.Statement.SetColumn(column, time.Now().Format(variable.DateFormat))
			}
		} else if destValueOf.Type().Kind() == reflect.Map {
			if b, column := structHasSpecialField("created_at", gormDB.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(variable.DateFormat)))
			}
			if b, column := structHasSpecialField("updated_at", gormDB.Statement.Dest); b {
				destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(variable.DateFormat)))
			}
		}
	}
}

func UpdateBeforeHook(gormDB *gorm.DB) {
	if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Struct {

		variable.ZapLog.Warn(my_errors.ErrorsGormDBUpdateParamsNotPtr)
	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Map {

	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(gormDB.Statement.Dest).Elem().Kind() == reflect.Struct {

		if b, column := structHasSpecialField("UpdatedAt", gormDB.Statement.Dest); b {
			gormDB.Statement.SetColumn(column, time.Now().Format(variable.DateFormat))
		}
	} else if reflect.TypeOf(gormDB.Statement.Dest).Kind() == reflect.Ptr && reflect.ValueOf(gormDB.Statement.Dest).Elem().Kind() == reflect.Map {
		if b, column := structHasSpecialField("updated_at", gormDB.Statement.Dest); b {
			destValueOf := reflect.ValueOf(gormDB.Statement.Dest).Elem()
			destValueOf.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(time.Now().Format(variable.DateFormat)))
		}
	}
}

func structHasSpecialField(fieldName string, anyStructPtr interface{}) (bool, string) {
	var tmp reflect.Type
	if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Map {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	} else if reflect.TypeOf(anyStructPtr).Kind() == reflect.Ptr && reflect.ValueOf(anyStructPtr).Elem().Kind() == reflect.Struct {
		destValueOf := reflect.ValueOf(anyStructPtr).Elem()
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Struct {

		destValueOf := anyStructPtr.(reflect.Value)
		tf := destValueOf.Type()
		for i := 0; i < tf.NumField(); i++ {
			if !tf.Field(i).Anonymous && tf.Field(i).Type.Kind() != reflect.Struct {
				if tf.Field(i).Name == fieldName {
					return true, getColumnNameFromGormTag(fieldName, tf.Field(i).Tag.Get("gorm"))
				}
			} else if tf.Field(i).Type.Kind() == reflect.Struct {
				tmp = tf.Field(i).Type
				for j := 0; j < tmp.NumField(); j++ {
					if tmp.Field(j).Name == fieldName {
						return true, getColumnNameFromGormTag(fieldName, tmp.Field(j).Tag.Get("gorm"))
					}
				}
			}
		}
	} else if reflect.Indirect(anyStructPtr.(reflect.Value)).Type().Kind() == reflect.Map {
		destValueOf := anyStructPtr.(reflect.Value)
		for _, item := range destValueOf.MapKeys() {
			if item.String() == fieldName {
				return true, fieldName
			}
		}
	}
	return false, ""
}

func getColumnNameFromGormTag(defaultColumn, TagValue string) (str string) {
	pos1 := strings.Index(TagValue, "column:")
	if pos1 == -1 {
		str = defaultColumn
		return
	} else {
		TagValue = TagValue[pos1+7:]
	}
	pos2 := strings.Index(TagValue, ";")
	if pos2 == -1 {
		str = TagValue
	} else {
		str = TagValue[:pos2]
	}
	return strings.ReplaceAll(str, " ", "")
}
