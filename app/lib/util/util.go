package util

import (
	"github.com/gorilla/securecookie"

	"github.com/mamoroom/echo-mvc/app/config"

	"encoding/base32"
	_ "fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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
	//fmt.Printf("[util] Generated Base32: %v\n", s)
	return s
}

func GetTimestampWithRand8(now time.Time) string {
	return CastInt64ToStr(now.Unix()) + "_" + GenerateBase32RandomKey(8)
}

func GetOauthProviderName(prefix string, suffix string) string {
	return prefix + "-" + suffix
}

func GetOauthProviderPrefix(oauth_provider_name string) string {
	parts := strings.Split(oauth_provider_name, "-")
	return parts[0]
}

func GetImageSuffixFromContentType(image_content_type string) string {
	parts := strings.Split(image_content_type, "/")
	return parts[1]
}

func GetFileNameSuffix(file_name string) string {
	parts := strings.Split(file_name, ".")
	return parts[len(parts)-1]
}

func GetTimeRfc3339(t time.Time) time.Time {
	_t, err := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	return _t
}

func GetBaseUrl() string {
	return conf.Server.Domain
}

func GetPortStr() string {
	return ":" + conf.Server.Internal.Port
}

func CastStrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func CastStrToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func CastStrToInt8(str string) (int8, error) {
	value, err := CastStrToInt64(str)
	return int8(value), err
}

func CastStrToInt(str string) (int, error) {
	value, err := CastStrToInt64(str)
	return int(value), err
}

func CastInt64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func CastUint64ToStr(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func CastUintToStr(u uint) string {
	return CastUint64ToStr(uint64(u))
}

func CastIntToStr(i int) string {
	return strconv.Itoa(i)
}

func CheckLoginAuthProviderValidation(str string) bool {
	return loginAuthProviderRegexp.Copy().Match([]byte(str))
	//return regexp.MustCompile(reg).Match([]byte(str))
}

func CheckChatAuthProviderValidation(str string) bool {
	return chatAuthProviderRegexp.Copy().Match([]byte(str))
	//return regexp.MustCompile(reg).Match([]byte(str))
}

func CheckImageContentTypeValidation(str string) bool {
	return imageContentTypeRegexp.Copy().Match([]byte(str))
	//return regexp.MustCompile(reg).Match([]byte(str))
}

func CheckDuplicateGameCodeValidation(err_msg string) bool {
	return duplicateGameCodeRegexp.Copy().Match([]byte(err_msg))
}

func CheckCralwerUaValidation(ua string) bool {
	return crawlerUaRegexp.Copy().Match([]byte(strings.ToLower(ua)))
}

func CheckCarrierDomainValidation(email string) bool {
	return carrierDomainRegexp.Copy().Match([]byte(email))
}

func GetRandInt(threshold int) int {
	rand.Seed(time.Now().UnixNano())
	// [1, threshold]で抽選
	return rand.Intn(threshold + 1)
}
