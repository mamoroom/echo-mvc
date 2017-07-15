package validate_data

import (
	"github.com/mamoroom/echo-mvc/app/lib/custom_io"
	"strings"
)

func IsInvalidNgWord(str string) bool {
	json_file_name := "ng"
	r := custom_io.GetRegexp(json_file_name, "validate", json_file_name, custom_io.GetValidateData)
	return r.Match([]byte(strings.ToLower(str)))
}
