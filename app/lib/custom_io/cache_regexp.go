package custom_io

import (
	"github.com/patrickmn/go-cache"

	"regexp"
	"strings"
)

func GetRegexp(key_prefix string, key_suffix string, cb_param interface{}, get_cb func(interface{}) []string) *regexp.Regexp {
	_r, found := data_cache.Get(cache_key(key_prefix, key_suffix))
	if !found {
		validate_data := get_cb(cb_param)
		regexp := regexp.MustCompile(strings.Join(validate_data, "|"))
		data_cache.Set(cache_key(key_prefix, key_suffix), &regexp, cache.DefaultExpiration)
		_r = &regexp
	}

	if _r == nil {
		panic("Target data or key error | [cache_key]" + cache_key(key_prefix, key_suffix))
	}
	_regexp, ok := _r.(**regexp.Regexp)
	if !ok {
		panic("Data format error | [cache_key]" + cache_key(key_prefix, key_suffix))
	}
	return (*_regexp).Copy()
}
