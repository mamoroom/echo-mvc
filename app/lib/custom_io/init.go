package custom_io

import (
	"github.com/mamoroom/echo-mvc/app/config"

	"github.com/patrickmn/go-cache"

	"encoding/gob"
	"time"
)

var conf = config.Conf
var data_cache *cache.Cache

func init() {
	data_cache = cache.New(time.Duration(conf.Data.ExpiresDuration)*time.Second, time.Duration(conf.Data.CleanUpDuration)*time.Second)

	// cache取り出し用: mapのdeep copyエンコーダーを作る
	gob.Register(map[string]interface{}{})
}
