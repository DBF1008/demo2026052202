package data_transfer

import (
	"encoding/json"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/http/validator/core/interf"
	"github.com/gin-gonic/gin"
	"time"
)

func DataAddContext(validatorInterface interf.ValidatorInterface, extraAddDataPrefix string, context *gin.Context) *gin.Context {
	var tempJson interface{}
	if tmpBytes, err1 := json.Marshal(validatorInterface); err1 == nil {
		if err2 := json.Unmarshal(tmpBytes, &tempJson); err2 == nil {
			if value, ok := tempJson.(map[string]interface{}); ok {
				for k, v := range value {
					context.Set(extraAddDataPrefix+k, v)
				}

				curDateTime := time.Now().Format(variable.DateFormat)
				context.Set(extraAddDataPrefix+"created_at", curDateTime)
				context.Set(extraAddDataPrefix+"updated_at", curDateTime)
				context.Set(extraAddDataPrefix+"deleted_at", curDateTime)
				return context
			}
		}
	}
	return nil
}
