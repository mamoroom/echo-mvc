package custom_io

import (
	"github.com/patrickmn/go-cache"

	"github.com/mamoroom/echo-mvc/app/lib/util"

	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetDataByKey(f_name string, key string) map[string]interface{} {
	_target, found := data_cache.Get(cache_key(f_name, key))
	if !found {
		all_data := get_all_data_from_file(f_name)
		all_data_map, ok := all_data.(map[string]interface{})
		if !ok {
			panic("Data format error | [file]" + f_name + ".json, [json_key]" + key)
		}
		for json_key, data := range all_data_map {
			d, ok := data.(map[string]interface{})
			if !ok {
				panic("Data format error | [file]" + f_name + ".json, [json_key]" + json_key)
			}
			data_cache.Set(cache_key(f_name, json_key), &d, cache.DefaultExpiration)
			if json_key == key {
				_target = &d
			}
		}
	}

	if _target == nil {
		panic("Target file or key error | [target_file]" + f_name + ".json, [target_key]" + key)
	}

	t, ok := _target.(*map[string]interface{})
	if !ok {
		panic("Data format error | [file]" + f_name + ".json, [json_key]" + key)
	}

	copy_t, err := deep_copy(*t)
	if err != nil {
		panic("Deep copy cache json error | [file]" + f_name + ".json, [json_key]" + key)
	}
	return copy_t
}

// callback for GetRegexp
func GetValidateData(key interface{}) []string {
	f_name := key.(string)
	d := get_all_data_from_file(f_name)
	data, ok := d.([]interface{})
	if !ok {
		panic("Data format error | [file]" + f_name + ".json")
	}

	validate_data := make([]string, len(data))
	for i, v := range data {
		// ``意味ない...
		validate_data[i] = `` + v.(string)
	}

	return validate_data
}

func get_all_data_from_file(f_name string) interface{} {
	f := util.GetRootPath() + "/config/data/" + f_name + ".json"
	fmt.Println("[io:READ] " + f)
	raw, err := ioutil.ReadFile(f)
	if err != nil {
		panic("Read invalid file | " + f)
	}
	var data interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		panic("Reading file is not json | " + f)
	}

	return data
}

func cache_key(key_prefix string, key_suffix string) string {
	return key_prefix + "__" + key_suffix
}

func deep_copy(m map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy map[string]interface{}
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}
