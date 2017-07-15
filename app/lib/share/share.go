package share

import (
	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/lib/util"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/url"
	"strings"
)

var conf = config.Conf
var aes_cipher cipher.Block

type Param struct {
	Id   int64  `json:"id"`
	Rank int64  `json:"rank"`
	Lang string `json:"lang"`
}

func init() {
	var err error
	aes_cipher, err = aes.NewCipher([]byte(conf.App.Share.SigningKey))
	if err != nil {
		panic("Invalid secret key is used for share crypt | " + err.Error())
	}
}

func CryptParam2Hash(param *Param) (string, error) {
	plaintext := param2text(param)
	code := make([]byte, aes.BlockSize+len(plaintext))
	iv := code[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(aes_cipher, iv)
	cfb.XORKeyStream(code[aes.BlockSize:], plaintext)

	return escapeString(code), nil
}

func escapeString(code []byte) string {
	base64_code := base64.StdEncoding.EncodeToString(code)
	base64_code_without_slash := strings.Replace(base64_code, "/", "-", -1)
	base64_code_without_slash_plus := strings.Replace(base64_code_without_slash, "+", "_", -1)
	return url.QueryEscape(base64_code_without_slash_plus)
}

func unEscapeString(hash string) ([]byte, error) {
	_hash_without_slash_plus, err := url.QueryUnescape(hash)
	if err != nil {
		return nil, err
	}
	_hash_without_plus := strings.Replace(_hash_without_slash_plus, "-", "/", -1)
	_hash := strings.Replace(_hash_without_plus, "_", "+", -1)
	code, err := base64.StdEncoding.DecodeString(_hash)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func DecryptHash2Param(hash string) (*Param, error) {
	code, err := unEscapeString(hash)
	if err != nil || len(code) < aes.BlockSize {
		return nil, err
	}
	iv := code[:aes.BlockSize]
	ciphertext := code[aes.BlockSize:]

	cfbdec := cipher.NewCFBDecrypter(aes_cipher, iv)
	cfbdec.XORKeyStream(ciphertext, ciphertext)

	param, err := text2param(string(ciphertext))
	if err != nil {
		return nil, err
	}
	return param, nil
}

func param2text(param *Param) []byte {
	str := util.CastInt64ToStr(param.Id) + ":" + util.CastInt64ToStr(param.Rank) + ":" + param.Lang
	return []byte(str)
}

func text2param(text string) (*Param, error) {
	parts := strings.Split(text, ":")
	id, err := util.CastStrToInt64(parts[0])
	if err != nil {
		return nil, err
	}
	rank, err := util.CastStrToInt64(parts[1])
	if err != nil {
		return nil, err
	}
	return &Param{
		Id:   id,
		Rank: rank,
		Lang: parts[2],
	}, nil
}
