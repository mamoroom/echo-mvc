package share_data

import (
	"github.com/mamoroom/echo-mvc/app/lib/custom_io"

	_ "fmt"
)

type ShareData struct {
	Data map[string]interface{}
}

func GetShareData(key string) (ShareData, error) {
	return ShareData{
		Data: custom_io.GetDataByKey("share", key),
	}, nil
}
