package util

import (
	"github.com/gorilla/securecookie"

	"github.com/mamoroom/echo-mvc/src/config"

	"encoding/base32"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var conf = config.Conf
var root_path = config.ROOT_PATH

func GetRootPath() string {
	return root_path
}

func GenerateRandomKey(length int) []byte {
	k := securecookie.GenerateRandomKey(length)
	if k == nil {
		panic("GenerateRandomKey Error")
	}
	return k
}

func GenerateBase32RandomKey(length int) string {
	s := strings.TrimRight(base32.StdEncoding.EncodeToString(GenerateRandomKey(length)), "=")
	fmt.Printf("[util] Generated Base32: %v\n", s)
	return s
}

func GetOauthProviderName(prefix string, suffix string) string {
	return prefix + "-" + suffix
}

func GetOauthProviderPrefix(oauth_provider_name string) string {
	parts := strings.Split(oauth_provider_name, "-")
	return parts[0]
}

func GetTimeRfc3339(t time.Time) time.Time {
	_t, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	return _t
}

func GetBaseUrl() string {
	return conf.Server.Domain + GetPortStr()
}

func GetPortStr() string {
	return ":" + conf.Server.Port
}

func CastStrToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func CastUint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// [todo]: 実行時間遅いらしい...
// https://developers.eure.jp/tech/golang-regexp/
func CheckRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
